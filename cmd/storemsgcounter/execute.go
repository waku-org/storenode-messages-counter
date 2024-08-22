package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/exp/maps"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/waku-org/go-waku/waku/v2/dnsdisc"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"github.com/waku-org/storenode-messages/internal/logging"
	"github.com/waku-org/storenode-messages/internal/metrics"
	"github.com/waku-org/storenode-messages/internal/persistence"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
const maxMsgHashesPerRequest = 100

type Application struct {
	node    *node.WakuNode
	metrics metrics.Metrics
	db      *persistence.DBStore
}

func Execute(ctx context.Context, options Options) error {
	// Set encoding for logs (console, json, ...)
	// Note that libp2p reads the encoding from GOLOG_LOG_FMT env var.
	logging.InitLogger(options.LogEncoding, options.LogOutput)

	logger := logging.Logger()

	logger.Warn("AppStart")

	var metricsServer *metrics.Server
	if options.EnableMetrics {
		metricsServer = metrics.NewMetricsServer(options.MetricsAddress, options.MetricsPort, logger)
		go metricsServer.Start()
	}

	var db *sql.DB
	var migrationFn func(*sql.DB, *zap.Logger) error
	db, migrationFn, err := persistence.ParseURL(options.DatabaseURL, logger)
	if err != nil {
		return err
	}

	dbStore, err := persistence.NewDBStore(
		options.ClusterID,
		options.FleetName,
		logger,
		persistence.WithDB(db),
		persistence.WithMigrations(migrationFn),
		persistence.WithRetentionPolicy(options.RetentionPolicy),
	)
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

	hostAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", options.Address, options.Port))
	if err != nil {
		return err
	}

	lvl, err := zapcore.ParseLevel(options.LogLevel)
	if err != nil {
		return err
	}

	wakuNode, err := node.New(
		node.WithLogLevel(lvl),
		node.WithNTP(),
		node.WithClusterID(uint16(options.ClusterID)),
		node.WithHostAddress(hostAddr),
	)
	if err != nil {
		return err
	}

	metrics := metrics.NewMetrics(options.ClusterID, options.FleetName, prometheus.DefaultRegisterer, logger)

	err = wakuNode.Start(ctx)
	if err != nil {
		return err
	}
	defer wakuNode.Stop()

	var storenodeIDs peer.IDSlice
	for _, s := range storenodes {
		wakuNode.Host().Peerstore().AddAddrs(s.ID, s.Addrs, peerstore.PermanentAddrTTL)
		storenodeIDs = append(storenodeIDs, s.ID)
	}

	err = dbStore.Start(ctx, wakuNode.Timesource())
	if err != nil {
		return err
	}

	application := &Application{
		node:    wakuNode,
		metrics: metrics,
		db:      dbStore,
	}

	go func() {
		missingMessagesTimer := time.NewTimer(0)
		defer missingMessagesTimer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-missingMessagesTimer.C:
				tmpUUID := uuid.New()
				runId := hex.EncodeToString(tmpUUID[:])
				runIdLogger := logger.With(zap.String("runId", runId))

				runIdLogger.Info("verifying message history...")
				shouldResetTimer, err := application.verifyHistory(ctx, runId, storenodeIDs, runIdLogger)
				if err != nil {
					runIdLogger.Error("could not verify message history", zap.Error(err))
				}
				runIdLogger.Info("verification complete")

				if shouldResetTimer {
					missingMessagesTimer.Reset(0)
				} else {
					missingMessagesTimer.Reset(timeInterval)
				}
			}
		}
	}()

	go func() {

		syncCheckTimer := time.NewTimer(0)
		defer syncCheckTimer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-syncCheckTimer.C:
				go func() {
					tmpUUID := uuid.New()
					runId := hex.EncodeToString(tmpUUID[:])
					runIdLogger := logger.With(zap.String("syncRunId", runId))

					err := application.checkMissingMessageStatus(ctx, storenodeIDs, runId, runIdLogger)
					if err != nil {
						logger.Error("could not recheck the status of missing messages", zap.Error(err))
						return
					}

					err = application.countMissingMessages(storenodeIDs)
					if err != nil {
						logger.Error("could not count missing messages", zap.Error(err))
						return
					}

					runIdLogger.Info("missing messages recheck complete")

					syncCheckTimer.Reset(30 * time.Minute)
				}()
			}
		}

	}()

	<-ctx.Done()

	return nil
}

var msgMapLock sync.Mutex
var msgMap map[pb.MessageHash]map[peer.ID]MessageExistence
var msgPubsubTopic map[pb.MessageHash]string

func (app *Application) verifyHistory(ctx context.Context, runId string, storenodes peer.IDSlice, logger *zap.Logger) (shouldResetTimer bool, err error) {

	// [MessageHash][StoreNode] = exists?
	msgMapLock.Lock()
	msgMap = make(map[pb.MessageHash]map[peer.ID]MessageExistence)
	msgPubsubTopic = make(map[pb.MessageHash]string)
	msgMapLock.Unlock()

	topicSyncStatus, err := app.db.GetTopicSyncStatus(ctx, options.PubSubTopics.Value())
	if err != nil {
		return false, err
	}

	tx, err := app.db.GetTrx(ctx)
	if err != nil {
		return false, err
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

		logger = logger.With(zap.String("topic", topic), zap.Timep("lastSyncTimestamp", lastSyncTimestamp))
		if lastSyncTimestamp != nil {
			app.metrics.RecordLastSyncDate(topic, *lastSyncTimestamp)
		}

		now := app.node.Timesource().Now()

		// Query is done with a delay
		startTime := now.Add(-(timeInterval + delay))
		if lastSyncTimestamp != nil {
			startTime = *lastSyncTimestamp
		}
		endTime := now.Add(-delay)

		if startTime.After(endTime) {
			logger.Warn("too soon to retrieve messages for topic")
			continue
		}

		// This avoids extremely large resultsets
		if endTime.Sub(startTime) > timeInterval {
			endTime = startTime.Add(timeInterval)
			shouldResetTimer = true
		}

		go func(topic string, startTime time.Time, endTime time.Time, logger *zap.Logger) {
			defer wg.Done()
			app.retrieveHistory(ctx, runId, storenodes, topic, startTime, endTime, tx, logger)
		}(topic, startTime, endTime, logger)
	}
	wg.Wait()

	// Verify for each storenode which messages are not available, and query
	// for their existence using message hash
	// ========================================================================
	msgsToVerify := make(map[peer.ID][]pb.MessageHash) // storenode -> msgHash
	msgMapLock.Lock()
	for msgHash, nodes := range msgMap {
		for _, s := range storenodes {
			if nodes[s] != Exists {
				msgsToVerify[s] = append(msgsToVerify[s], msgHash)
			}
		}
	}
	msgMapLock.Unlock()

	wg = sync.WaitGroup{}
	for peerID, messageHashes := range msgsToVerify {
		wg.Add(1)
		go func(peerID peer.ID, messageHashes []pb.MessageHash) {
			defer wg.Done()

			onResult := func(result *store.Result) {
				msgMapLock.Lock()
				for _, mkv := range result.Messages() {
					hash := mkv.WakuMessageHash()
					_, ok := msgMap[hash]
					if !ok {
						msgMap[hash] = make(map[peer.ID]MessageExistence)
					}
					msgMap[hash][result.PeerID()] = Exists
				}

				for _, msgHash := range messageHashes {
					if msgMap[msgHash][result.PeerID()] != Exists {
						msgMap[msgHash][result.PeerID()] = DoesNotExist
					}
				}
				msgMapLock.Unlock()
			}

			app.verifyMessageExistence(ctx, runId, peerID, messageHashes, onResult, logger)
		}(peerID, messageHashes)
	}
	wg.Wait()

	// If a message is not available, store in DB in which store nodes it wasnt
	// available
	// ========================================================================
	msgMapLock.Lock()
	defer msgMapLock.Unlock()

	missingInSummary := make(map[peer.ID]int)
	unknownInSummary := make(map[peer.ID]int)
	totalMissingMessages := 0

	for msgHash, nodes := range msgMap {
		var missingIn []peer.ID
		var unknownIn []peer.ID

		for _, s := range storenodes {
			if nodes[s] == DoesNotExist {
				missingIn = append(missingIn, s)
				missingInSummary[s]++
			} else if nodes[s] == Unknown {
				unknownIn = append(unknownIn, s)
				unknownInSummary[s]++
			}
		}

		if len(missingIn) != 0 {
			logger.Info("missing message identified", zap.Stringer("hash", msgHash), zap.String("pubsubTopic", msgPubsubTopic[msgHash]), zap.Int("num_nodes", len(missingIn)))
			err := app.db.RecordMessage(runId, tx, msgHash, msgPubsubTopic[msgHash], missingIn, "does_not_exist")
			if err != nil {
				return false, err
			}
			totalMissingMessages++
		}

		if len(unknownIn) != 0 {
			logger.Info("message with unknown state identified", zap.Stringer("hash", msgHash), zap.String("pubsubTopic", msgPubsubTopic[msgHash]), zap.Int("num_nodes", len(missingIn)))
			err = app.db.RecordMessage(runId, tx, msgHash, msgPubsubTopic[msgHash], unknownIn, "unknown")
			if err != nil {
				return false, err
			}
		}
	}

	for _, s := range storenodes {
		missingCnt := missingInSummary[s]
		app.metrics.RecordMissingMessages(s, "does_not_exist", missingCnt)
		logger.Info("missing message summary", zap.Stringer("storenode", s), zap.Int("numMsgs", missingCnt))

		unknownCnt := unknownInSummary[s]
		app.metrics.RecordMissingMessages(s, "unknown", unknownCnt)
		logger.Info("messages that could not be verified summary", zap.Stringer("storenode", s), zap.Int("numMsgs", missingCnt))
	}

	logger.Info("total missing messages", zap.Int("total", totalMissingMessages))
	app.metrics.RecordTotalMissingMessages(totalMissingMessages)

	return shouldResetTimer, nil
}

func (app *Application) checkMissingMessageStatus(ctx context.Context, storenodes []peer.ID, runId string, logger *zap.Logger) error {
	now := app.node.Timesource().Now()

	from := now.Add(-2 * time.Hour)
	to := now.Add(-time.Hour)

	logger.Info("rechecking missing messages status", zap.Time("from", from), zap.Time("to", to), zap.Uint("clusterID", options.ClusterID))

	// Get all messages whose status is missing or does not exist, and the column found_on_recheck is false
	// if found, set found_on_recheck to true
	missingMessages, err := app.db.GetMissingMessages(from, to)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, storenodeID := range storenodes {
		wg.Add(1)
		go func(peerID peer.ID, messageHashes []pb.MessageHash) {
			defer wg.Done()

			foundMissingMessages := make(map[pb.MessageHash]struct{})
			app.verifyMessageExistence(ctx, runId, peerID, messageHashes, func(result *store.Result) {
				for _, mkv := range result.Messages() {
					foundMissingMessages[mkv.WakuMessageHash()] = struct{}{}
				}
			}, logger)

			err := app.db.MarkMessagesAsFound(peerID, maps.Keys(foundMissingMessages), options.ClusterID)
			if err != nil {
				logger.Error("could not mark messages as found", zap.Error(err))
				return
			}

			cnt := len(messageHashes) - len(foundMissingMessages)
			app.metrics.RecordMissingMessagesPrevHour(peerID, cnt)
			logger.Info("missingMessages for the previous hour", zap.Stringer("storenode", peerID), zap.Int("cnt", cnt))

		}(storenodeID, missingMessages[storenodeID])
	}
	wg.Wait()

	return nil
}

func (app *Application) countMissingMessages(storenodes []peer.ID) error {

	// not including last two hours in now to let sync work
	now := app.node.Timesource().Now().Add(-2 * time.Hour)

	// Count messages in last day (not including last two hours)
	results, err := app.db.CountMissingMessages(now.Add(-24*time.Hour), now, options.ClusterID)
	if err != nil {
		return err
	}
	for storenode, cnt := range results {
		app.metrics.RecordMissingMessagesLastDay(storenode, cnt)
	}

	// Count messages in last week (not including last two hours)
	results, err = app.db.CountMissingMessages(now.Add(-24*time.Hour*7), now, options.ClusterID)
	if err != nil {
		return err
	}
	for _, storenodeID := range storenodes {
		app.metrics.RecordMissingMessagesLastWeek(storenodeID, results[storenodeID])
	}
	return nil
}

func (app *Application) fetchStoreNodeMessages(ctx context.Context, runId string, storenodeID peer.ID, topic string, startTime time.Time, endTime time.Time, logger *zap.Logger) {
	var result *store.Result
	var err error

	queryLogger := logger.With(zap.Stringer("storenode", storenodeID), zap.Int64("startTime", startTime.UnixNano()), zap.Int64("endTime", endTime.UnixNano()))

	retry := true
	success := false
	count := 1
	for retry && count <= maxAttempts {
		queryLogger.Info("retrieving message history for topic!", zap.Int("attempt", count))

		tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		result, err = app.node.Store().Query(tCtx, store.FilterCriteria{
			ContentFilter: protocol.NewContentFilter(topic),
			TimeStart:     proto.Int64(startTime.UnixNano()),
			TimeEnd:       proto.Int64(endTime.UnixNano()),
		}, store.WithPeer(storenodeID), store.WithPaging(false, 100), store.IncludeData(false))
		cancel()

		if err != nil {
			queryLogger.Error("could not query storenode", zap.Error(err), zap.Int("attempt", count))
			time.Sleep(2 * time.Second)
		} else {
			queryLogger.Info("messages available", zap.Int("len", len(result.Messages())))
			retry = false
			success = true
		}
		count++
	}

	if !success {
		queryLogger.Error("storenode not available")
		err := app.db.RecordStorenodeUnavailable(runId, storenodeID)
		if err != nil {
			queryLogger.Error("could not store node unavailable", zap.Error(err))
		}
		app.metrics.RecordStorenodeAvailability(storenodeID, false)
		return
	}

	app.metrics.RecordStorenodeAvailability(storenodeID, true)

	for !result.IsComplete() {
		msgMapLock.Lock()
		for _, mkv := range result.Messages() {
			hash := mkv.WakuMessageHash()
			_, ok := msgMap[hash]
			if !ok {
				msgMap[hash] = make(map[peer.ID]MessageExistence)
			}
			msgMap[hash][storenodeID] = Exists
			msgPubsubTopic[hash] = topic
		}
		msgMapLock.Unlock()

		retry := true
		success := false
		count := 1
		cursorLogger := queryLogger.With(zap.String("cursor", hex.EncodeToString(result.Cursor())))
		for retry && count <= maxAttempts {
			cursorLogger.Info("retrieving next page")
			tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
			err = result.Next(tCtx)
			cancel()
			if err != nil {
				cursorLogger.Error("could not query storenode", zap.Error(err))
				time.Sleep(2 * time.Second)
			} else {
				cursorLogger.Info("more messages available", zap.Int("len", len(result.Messages())))
				retry = false
				success = true
			}
			count++
		}

		if !success {
			cursorLogger.Error("storenode not available")
			err := app.db.RecordStorenodeUnavailable(runId, storenodeID)
			if err != nil {
				cursorLogger.Error("could not store recordnode unavailable", zap.Error(err))
			}
			app.metrics.RecordStorenodeAvailability(storenodeID, false)
			return
		}

		app.metrics.RecordStorenodeAvailability(storenodeID, true)
	}
}

func (app *Application) retrieveHistory(ctx context.Context, runId string, storenodes peer.IDSlice, topic string, startTime time.Time, endTime time.Time, tx *sql.Tx, logger *zap.Logger) {
	// Determine if the messages exist across all nodes
	wg := sync.WaitGroup{}
	for _, storePeerID := range storenodes {
		wg.Add(1)
		go func(peerID peer.ID) {
			defer wg.Done()
			app.fetchStoreNodeMessages(ctx, runId, peerID, topic, startTime, endTime, logger)
		}(storePeerID)
	}

	wg.Wait()

	// Update db with last sync time
	err := app.db.UpdateTopicSyncState(tx, topic, endTime)
	if err != nil {
		logger.Panic("could not update topic sync state", zap.Error(err))
	}

	app.metrics.RecordLastSyncDate(topic, endTime)

}

func (app *Application) verifyMessageExistence(ctx context.Context, runId string, peerID peer.ID, messageHashes []pb.MessageHash, onResult func(result *store.Result), logger *zap.Logger) {
	if len(messageHashes) == 0 {
		return
	}

	peerInfo := app.node.Host().Peerstore().PeerInfo(peerID)

	queryLogger := logger.With(zap.Stringer("storenode", peerID))

	wg := sync.WaitGroup{}
	// Split into batches
	for i := 0; i < len(messageHashes); i += maxMsgHashesPerRequest {
		j := i + maxMsgHashesPerRequest
		if j > len(messageHashes) {
			j = len(messageHashes)
		}

		wg.Add(1)
		go func(messageHashes []pb.MessageHash) {
			defer wg.Done()

			var result *store.Result
			var err error

			retry := true
			success := false
			count := 1
			for retry && count <= maxAttempts {
				queryLogger.Info("querying by hash", zap.Int("attempt", count))
				tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
				result, err = app.node.Store().QueryByHash(tCtx, messageHashes, store.IncludeData(false), store.WithPeer(peerInfo.ID), store.WithPaging(false, 100))
				cancel()
				if err != nil {
					queryLogger.Error("could not query storenode", zap.Error(err), zap.Int("attempt", count))
					time.Sleep(2 * time.Second)
				} else {
					queryLogger.Info("hashes available", zap.Int("len", len(result.Messages())))
					retry = false
					success = true
				}
				count++
			}

			if !success {
				queryLogger.Error("storenode not available")
				err := app.db.RecordStorenodeUnavailable(runId, peerID)
				if err != nil {
					queryLogger.Error("could not store recordnode unavailable", zap.Error(err))
				}
				app.metrics.RecordStorenodeAvailability(peerID, false)
				return
			}

			app.metrics.RecordStorenodeAvailability(peerID, true)

			for !result.IsComplete() {
				onResult(result)

				retry := true
				success := false
				count := 1
				for retry && count <= maxAttempts {
					queryLogger.Info("executing next while querying hashes", zap.String("cursor", hexutil.Encode(result.Cursor())), zap.Int("attempt", count))
					tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
					err = result.Next(tCtx)
					cancel()
					if err != nil {
						queryLogger.Error("could not query storenode", zap.String("cursor", hexutil.Encode(result.Cursor())), zap.Error(err), zap.Int("attempt", count))
						time.Sleep(2 * time.Second)
					} else {
						queryLogger.Info("more hashes available", zap.String("cursor", hex.EncodeToString(result.Cursor())), zap.Int("len", len(result.Messages())))
						retry = false
						success = true
					}
					count++
				}

				if !success {
					queryLogger.Error("storenode not available", zap.String("cursor", hexutil.Encode(result.Cursor())))
					err := app.db.RecordStorenodeUnavailable(runId, peerID)
					if err != nil {
						logger.Error("could not store recordnode unavailable", zap.Error(err), zap.String("cursor", hex.EncodeToString(result.Cursor())), zap.Stringer("storenode", peerInfo))
					}
					app.metrics.RecordStorenodeAvailability(peerID, false)
					return
				}

				app.metrics.RecordStorenodeAvailability(peerID, true)
			}
		}(messageHashes[i:j])
	}

	wg.Wait()
}
