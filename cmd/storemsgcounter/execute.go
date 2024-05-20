package main

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/multiformats/go-multiaddr"
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
			err := verifyHistory(ctx, wakuNode, dbStore, logger)
			if err != nil {
				return err
			}
			logger.Info("verification complete")

			timer.Reset(timeInterval)
		}
	}
}

var msgMapLock sync.Mutex
var msgMap map[pb.MessageHash]map[string]MessageExistence
var msgAttr map[pb.MessageHash]MessageAttr

func verifyHistory(ctx context.Context, wakuNode *node.WakuNode, dbStore *persistence.DBStore, logger *zap.Logger) error {
	runId := uuid.New().String()

	logger = logger.With(zap.String("runId", runId))

	// [MessageHash][StoreNode] = exists?
	msgMapLock.Lock()
	msgMap = make(map[pb.MessageHash]map[string]MessageExistence)
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
			retrieveHistory(ctx, runId, topic, lastSyncTimestamp, wakuNode, dbStore, tx, logger)
		}(topic, lastSyncTimestamp)
	}
	wg.Wait()

	// Verify for each storenode which messages are not available, and query
	// for their existence using message hash
	// ========================================================================
	msgsToVerify := make(map[string][]pb.MessageHash) // storenode -> msgHash
	msgMapLock.Lock()
	for msgHash, nodes := range msgMap {
		for _, node := range options.StoreNodes {
			nodeStr := node.String()
			if nodes[nodeStr] != Exists {
				msgsToVerify[nodeStr] = append(msgsToVerify[nodeStr], msgHash)
			}
		}
	}
	msgMapLock.Unlock()

	wg = sync.WaitGroup{}
	for node, messageHashes := range msgsToVerify {
		wg.Add(1)
		go func(node string, messageHashes []pb.MessageHash) {
			defer wg.Done()
			nodeMultiaddr, _ := multiaddr.NewMultiaddr(node)
			verifyMessageExistence(ctx, runId, tx, nodeMultiaddr, messageHashes, wakuNode, dbStore, logger)
		}(node, messageHashes)
	}
	wg.Wait()

	// If a message is not available, store in DB in which store nodes it wasnt
	// available and its timestamp
	// ========================================================================
	msgMapLock.Lock()
	defer msgMapLock.Unlock()
	for msgHash, nodes := range msgMap {
		var missingIn []string
		var unknownIn []string
		for _, node := range options.StoreNodes {
			nodeStr := node.String()
			if nodes[nodeStr] == DoesNotExist {
				missingIn = append(missingIn, nodeStr)
			} else if nodes[nodeStr] == Unknown {
				unknownIn = append(unknownIn, nodeStr)
			}
		}

		err := dbStore.RecordMessage(runId, tx, msgHash, options.ClusterID, msgAttr[msgHash].PubsubTopic, msgAttr[msgHash].Timestamp, missingIn, "does_not_exist")
		if err != nil {
			return err
		}

		err = dbStore.RecordMessage(runId, tx, msgHash, options.ClusterID, msgAttr[msgHash].PubsubTopic, msgAttr[msgHash].Timestamp, unknownIn, "unknown")
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func retrieveHistory(ctx context.Context, runId string, topic string, lastSyncTimestamp *time.Time, wakuNode *node.WakuNode, dbStore *persistence.DBStore, tx *sql.Tx, logger *zap.Logger) {
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
	for _, node := range options.StoreNodes {
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
			}, store.WithPeerAddr(node))
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
			err := dbStore.RecordStorenodeUnavailable(runId, node.String())
			if err != nil {
				logger.Error("could not store recordnode unavailable", zap.Error(err), zap.Stringer("storenode", node))
			}

		} else {

		iteratorLbl:
			for !result.IsComplete() {
				msgMapLock.Lock()
				for _, mkv := range result.Messages() {
					hash := mkv.WakuMessageHash()
					_, ok := msgMap[hash]
					if !ok {
						msgMap[hash] = make(map[string]MessageExistence)
					}
					msgMap[hash][node.String()] = Exists
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
					// TODO: Notify that storenode was not available from X to Y time
					logger.Error("storenode not available",
						zap.Stringer("storenode", node),
						zap.Time("startTime", startTime),
						zap.Time("endTime", endTime),
						zap.String("topic", topic),
						zap.String("cursor", hexutil.Encode(result.Cursor())))
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

func verifyMessageExistence(ctx context.Context, runId string, tx *sql.Tx, nodeAddr multiaddr.Multiaddr, messageHashes []pb.MessageHash, wakuNode *node.WakuNode, dbStore *persistence.DBStore, logger *zap.Logger) {
	storeNodeFailure := false
	var result *store.Result
	var err error

queryLbl:
	for i := 0; i < maxAttempts; i++ {
		result, err = wakuNode.Store().QueryByHash(ctx, messageHashes, store.IncludeData(false), store.WithPeerAddr(nodeAddr))
		if err != nil {
			logger.Error("could not query storenode", zap.Stringer("storenode", nodeAddr), zap.Error(err))
			storeNodeFailure = true
			time.Sleep(2 * time.Second)
		} else {
			storeNodeFailure = false
			break queryLbl
		}
	}

	if storeNodeFailure {
		logger.Error("storenode not available",
			zap.Stringer("storenode", nodeAddr),
			zap.Stringers("hashes", messageHashes))

		err := dbStore.RecordStorenodeUnavailable(runId, nodeAddr.String())
		if err != nil {
			logger.Error("could not store recordnode unavailable", zap.Error(err), zap.Stringer("storenode", nodeAddr))
		}
	} else {
		for !result.IsComplete() {
			nodeAddrStr := nodeAddr.String()

			msgMapLock.Lock()
			for _, mkv := range result.Messages() {
				hash := mkv.WakuMessageHash()
				_, ok := msgMap[hash]
				if !ok {
					msgMap[hash] = make(map[string]MessageExistence)
				}
				msgMap[hash][nodeAddrStr] = Exists
			}

			for _, msgHash := range messageHashes {
				if msgMap[msgHash][nodeAddrStr] != Exists {
					msgMap[msgHash][nodeAddrStr] = DoesNotExist
				}
			}

			msgMapLock.Unlock()

			storeNodeFailure = false

		nextRetryLbl:
			for i := 0; i < maxAttempts; i++ {
				err = result.Next(ctx)
				if err != nil {
					logger.Error("could not query storenode", zap.Stringer("storenode", nodeAddr), zap.Error(err))
					storeNodeFailure = true
					time.Sleep(2 * time.Second)
				} else {
					storeNodeFailure = false
					break nextRetryLbl
				}
			}

			if storeNodeFailure {
				logger.Error("storenode not available",
					zap.Stringer("storenode", nodeAddr),
					zap.Stringers("hashes", messageHashes),
					zap.String("cursor", hexutil.Encode(result.Cursor())))

				err := dbStore.RecordStorenodeUnavailable(runId, nodeAddr.String())
				if err != nil {
					logger.Error("could not store recordnode unavailable", zap.Error(err), zap.Stringer("storenode", nodeAddr))
				}
			}

		}
	}
}
