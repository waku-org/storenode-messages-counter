package main

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/waku-org/go-waku/waku/v2/dnsdisc"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"github.com/waku-org/storenode-messages/internal/logging"
	"github.com/waku-org/storenode-messages/internal/persistence"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type MessageExistence int

const (
	Unknown MessageExistence = iota
	Exists
	DoesNotExist
)

const timeInterval = 2 * time.Minute
const delay = 5 * time.Minute
const maxAttempts = 3

type MessageAttr struct {
	Timestamp   uint64
	PubsubTopic string
}

func Execute(ctx context.Context, options Options) error {
	// Set encoding for logs (console, json, ...)
	// Note that libp2p reads the encoding from GOLOG_LOG_FMT env var.
	logging.InitLogger(options.LogEncoding, options.LogOutput)

	logger := logging.Logger()

	var db *sql.DB
	var migrationFn func(*sql.DB, *zap.Logger) error
	db, migrationFn, err := persistence.ParseURL(options.DatabaseURL, logger)
	if err != nil {
		return err
	}

	dbStore, err := persistence.NewDBStore(logger, persistence.WithDB(db), persistence.WithMigrations(migrationFn), persistence.WithRetentionPolicy(options.RetentionPolicy))
	if err != nil {
		return err
	}
	defer dbStore.Stop()

	var discoveredNodes []dnsdisc.DiscoveredNode
	if len(options.DNSDiscoveryURLs.Value()) != 0 {
		discoveredNodes = node.GetNodesFromDNSDiscovery(logger, ctx, options.DNSDiscoveryNameserver, options.DNSDiscoveryURLs.Value())
	}

	var storenodes []peer.AddrInfo
	for _, node := range discoveredNodes {
		if len(node.PeerInfo.Addrs) == 0 {
			continue
		}
		storenodes = append(storenodes, node.PeerInfo)
	}

	for _, node := range options.StoreNodes {
		pInfo, err := peer.AddrInfosFromP2pAddrs(node)
		if err != nil {
			return err
		}

		storenodes = append(storenodes, pInfo...)
	}

	if len(storenodes) == 0 {
		return errors.New("no storenodes specified")
	}

	wakuNode, err := node.New(
		node.WithNTP(),
		node.WithClusterID(uint16(options.ClusterID)),
	)
	if err != nil {
		return err
	}
	err = wakuNode.Start(ctx)
	if err != nil {
		return err
	}
	defer wakuNode.Stop()

	for _, s := range storenodes {
		wakuNode.Host().Peerstore().AddAddrs(s.ID, s.Addrs, peerstore.PermanentAddrTTL)
	}

	err = dbStore.Start(ctx, wakuNode.Timesource())
	if err != nil {
		return err
	}

	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			logger.Info("verifying message history...")
			err := verifyHistory(ctx, storenodes, wakuNode, dbStore, logger)
			if err != nil {
				return err
			}
			logger.Info("verification complete")

			timer.Reset(timeInterval)
		}
	}
}

var msgMapLock sync.Mutex
var msgMap map[pb.MessageHash]map[peer.ID]MessageExistence
var msgAttr map[pb.MessageHash]MessageAttr

func verifyHistory(ctx context.Context, storenodes []peer.AddrInfo, wakuNode *node.WakuNode, dbStore *persistence.DBStore, logger *zap.Logger) error {
	runId := uuid.New().String()

	logger = logger.With(zap.String("runId", runId))

	// [MessageHash][StoreNode] = exists?
	msgMapLock.Lock()
	msgMap = make(map[pb.MessageHash]map[peer.ID]MessageExistence)
	msgAttr = make(map[pb.MessageHash]MessageAttr)
	msgMapLock.Unlock()

	topicSyncStatus, err := dbStore.GetTopicSyncStatus(ctx, options.ClusterID, options.PubSubTopics.Value())
	if err != nil {
		return err
	}

	tx, err := dbStore.GetTrx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
			return
		}
		// don't shadow original error
		_ = tx.Rollback()
	}()

	wg := sync.WaitGroup{}
	for topic, lastSyncTimestamp := range topicSyncStatus {
		wg.Add(1)
		go func(topic string, lastSyncTimestamp *time.Time) {
			defer wg.Done()
			retrieveHistory(ctx, runId, storenodes, topic, lastSyncTimestamp, wakuNode, dbStore, tx, logger)
		}(topic, lastSyncTimestamp)
	}
	wg.Wait()

	// Verify for each storenode which messages are not available, and query
	// for their existence using message hash
	// ========================================================================
	msgsToVerify := make(map[peer.ID][]pb.MessageHash) // storenode -> msgHash
	msgMapLock.Lock()
	for msgHash, nodes := range msgMap {
		for _, node := range storenodes {
			if nodes[node.ID] != Exists {
				msgsToVerify[node.ID] = append(msgsToVerify[node.ID], msgHash)
			}
		}
	}
	msgMapLock.Unlock()

	wg = sync.WaitGroup{}
	for peerID, messageHashes := range msgsToVerify {
		wg.Add(1)
		go func(peerID peer.ID, messageHashes []pb.MessageHash) {
			defer wg.Done()
			verifyMessageExistence(ctx, runId, tx, peerID, messageHashes, wakuNode, dbStore, logger)
		}(peerID, messageHashes)
	}
	wg.Wait()

	// If a message is not available, store in DB in which store nodes it wasnt
	// available and its timestamp
	// ========================================================================
	msgMapLock.Lock()
	defer msgMapLock.Unlock()

	for msgHash, nodes := range msgMap {
		var missingIn []peer.AddrInfo
		var unknownIn []peer.AddrInfo
		for _, node := range storenodes {
			if nodes[node.ID] == DoesNotExist {
				missingIn = append(missingIn, node)
			} else if nodes[node.ID] == Unknown {
				unknownIn = append(unknownIn, node)
			}
		}

		if len(missingIn) != 0 {
			logger.Info("missing message identified", zap.Stringer("hash", msgHash), zap.String("pubsubTopic", msgAttr[msgHash].PubsubTopic), zap.Int("num_nodes", len(missingIn)))
			err := dbStore.RecordMessage(runId, tx, msgHash, options.ClusterID, msgAttr[msgHash].PubsubTopic, msgAttr[msgHash].Timestamp, missingIn, "does_not_exist")
			if err != nil {
				return err
			}
		}

		if len(unknownIn) != 0 {
			logger.Debug("message with unknown state identified", zap.Stringer("hash", msgHash), zap.String("pubsubTopic", msgAttr[msgHash].PubsubTopic), zap.Int("num_nodes", len(missingIn)))
			err = dbStore.RecordMessage(runId, tx, msgHash, options.ClusterID, msgAttr[msgHash].PubsubTopic, msgAttr[msgHash].Timestamp, unknownIn, "unknown")
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func retrieveHistory(ctx context.Context, runId string, storenodes []peer.AddrInfo, topic string, lastSyncTimestamp *time.Time, wakuNode *node.WakuNode, dbStore *persistence.DBStore, tx *sql.Tx, logger *zap.Logger) {
	logger = logger.With(zap.String("topic", topic), zap.Timep("lastSyncTimestamp", lastSyncTimestamp))

	now := wakuNode.Timesource().Now()

	// Query is done with a delay
	startTime := now.Add(-(timeInterval + delay))
	if lastSyncTimestamp != nil {
		startTime = *lastSyncTimestamp
	}
	endTime := now.Add(-delay)

	if startTime.After(endTime) {
		logger.Warn("too soon to retrieve messages for topic")
		return
	}

	// Determine if the messages exist across all nodes
	for _, node := range storenodes {
		storeNodeFailure := false
		var result *store.Result
		var err error

		logger.Info("retrieving message history for topic", zap.Stringer("storenode", node), zap.Int64("from", startTime.UnixNano()), zap.Int64("to", endTime.UnixNano()))

	queryLbl:
		for i := 0; i < maxAttempts; i++ {
			result, err = wakuNode.Store().Query(ctx, store.FilterCriteria{
				ContentFilter: protocol.NewContentFilter(topic),
				TimeStart:     proto.Int64(startTime.UnixNano()),
				TimeEnd:       proto.Int64(endTime.UnixNano()),
			}, store.WithPeer(node.ID))
			if err != nil {
				logger.Error("could not query storenode", zap.Stringer("storenode", node), zap.Error(err))
				storeNodeFailure = true
				time.Sleep(2 * time.Second)
			} else {
				logger.Debug("messages available?", zap.Int("len", len(result.Messages())))
				storeNodeFailure = false
				break queryLbl
			}
		}

		if storeNodeFailure {
			logger.Error("storenode not available", zap.Stringer("storenode", node), zap.Time("startTime", startTime), zap.Time("endTime", endTime))
			err := dbStore.RecordStorenodeUnavailable(runId, node)
			if err != nil {
				logger.Error("could not store node unavailable", zap.Error(err), zap.Stringer("storenode", node))
			}

		} else {

		iteratorLbl:
			for !result.IsComplete() {
				msgMapLock.Lock()
				for _, mkv := range result.Messages() {
					hash := mkv.WakuMessageHash()
					_, ok := msgMap[hash]
					if !ok {
						msgMap[hash] = make(map[peer.ID]MessageExistence)
					}
					msgMap[hash][node.ID] = Exists
					msgAttr[hash] = MessageAttr{
						Timestamp:   uint64(mkv.Message.GetTimestamp()),
						PubsubTopic: mkv.GetPubsubTopic(),
					}
				}
				msgMapLock.Unlock()

				storeNodeFailure := false

			nextRetryLbl:
				for i := 0; i < maxAttempts; i++ {
					err = result.Next(ctx)
					if err != nil {
						logger.Error("could not query storenode", zap.Stringer("storenode", node), zap.Error(err))
						storeNodeFailure = true
						time.Sleep(2 * time.Second)
					} else {
						storeNodeFailure = false
						break nextRetryLbl
					}
				}

				if storeNodeFailure {
					logger.Error("storenode not available",
						zap.Stringer("storenode", node),
						zap.Time("startTime", startTime),
						zap.Time("endTime", endTime),
						zap.String("topic", topic),
						zap.String("cursor", hexutil.Encode(result.Cursor())))

					err := dbStore.RecordStorenodeUnavailable(runId, node)
					if err != nil {
						logger.Error("could not store recordnode unavailable", zap.Error(err), zap.Stringer("storenode", node))
					}
					break iteratorLbl
				}
			}
		}
	}

	// Update db with last sync time
	err := dbStore.UpdateTopicSyncState(tx, options.ClusterID, topic, endTime)
	if err != nil {
		logger.Panic("could not update topic sync state", zap.Error(err))
	}
}

func verifyMessageExistence(ctx context.Context, runId string, tx *sql.Tx, peerID peer.ID, messageHashes []pb.MessageHash, wakuNode *node.WakuNode, dbStore *persistence.DBStore, logger *zap.Logger) {
	storeNodeFailure := false
	var result *store.Result
	var err error

	peerInfo := wakuNode.Host().Peerstore().PeerInfo(peerID)

queryLbl:
	for i := 0; i < maxAttempts; i++ {
		result, err = wakuNode.Store().QueryByHash(ctx, messageHashes, store.IncludeData(false), store.WithPeer(peerInfo.ID))
		if err != nil {
			logger.Error("could not query storenode", zap.Stringer("storenode", peerInfo), zap.Error(err))
			storeNodeFailure = true
			time.Sleep(2 * time.Second)
		} else {
			storeNodeFailure = false
			break queryLbl
		}
	}

	if storeNodeFailure {
		logger.Error("storenode not available",
			zap.Stringer("storenode", peerInfo),
			zap.Stringers("hashes", messageHashes))

		err := dbStore.RecordStorenodeUnavailable(runId, peerInfo)
		if err != nil {
			logger.Error("could not store recordnode unavailable", zap.Error(err), zap.Stringer("storenode", peerInfo))
		}
	} else {
		for !result.IsComplete() {
			msgMapLock.Lock()
			for _, mkv := range result.Messages() {
				hash := mkv.WakuMessageHash()
				_, ok := msgMap[hash]
				if !ok {
					msgMap[hash] = make(map[peer.ID]MessageExistence)
				}
				msgMap[hash][peerInfo.ID] = Exists
			}

			for _, msgHash := range messageHashes {
				if msgMap[msgHash][peerInfo.ID] != Exists {
					msgMap[msgHash][peerInfo.ID] = DoesNotExist
				}
			}

			msgMapLock.Unlock()

			storeNodeFailure = false

		nextRetryLbl:
			for i := 0; i < maxAttempts; i++ {
				err = result.Next(ctx)
				if err != nil {
					logger.Error("could not query storenode", zap.Stringer("storenode", peerInfo), zap.Error(err))
					storeNodeFailure = true
					time.Sleep(2 * time.Second)
				} else {
					storeNodeFailure = false
					break nextRetryLbl
				}
			}

			if storeNodeFailure {
				logger.Error("storenode not available",
					zap.Stringer("storenode", peerInfo),
					zap.Stringers("hashes", messageHashes),
					zap.String("cursor", hexutil.Encode(result.Cursor())))

				err := dbStore.RecordStorenodeUnavailable(runId, peerInfo)
				if err != nil {
					logger.Error("could not store recordnode unavailable", zap.Error(err), zap.Stringer("storenode", peerInfo))
				}
			}

		}
	}
}
