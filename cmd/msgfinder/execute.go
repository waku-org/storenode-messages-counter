package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
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

const maxAttempts = 3

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

	dbStore, err := persistence.NewDBStore(logger, persistence.WithDB(db), persistence.WithMigrations(migrationFn))
	if err != nil {
		return err
	}
	defer dbStore.Stop()

	wakuNode, err := node.New(
		node.WithNTP(),
		node.WithClusterID(uint16(options.ClusterID)),
	)
	err = wakuNode.Start(ctx)
	if err != nil {
		return err
	}
	defer wakuNode.Stop()

	err = dbStore.Start(ctx, wakuNode.Timesource())
	if err != nil {
		return err
	}

	timeInterval := 2 * time.Minute
	delay := 5 * time.Minute

	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			// [MessageHash][StoreNode] = exists?
			msgMap := make(map[pb.MessageHash]map[string]MessageExistence)
			msgTopic := make(map[pb.MessageHash]string)

			topicSyncStatus, err := dbStore.GetTopicSyncStatus(ctx, options.ClusterID, options.PubSubTopics.Value())
			if err != nil {
				return err
			}

			trx, err := dbStore.GetTrx(ctx)
			if err != nil {
				return err
			}

			// TODO: commit or revert trx

			for topic, lastSyncTimestamp := range topicSyncStatus {
				now := wakuNode.Timesource().Now()

				// Query is done with a delay
				startTime := now.Add(-(timeInterval + delay))
				if lastSyncTimestamp != nil {
					startTime = *lastSyncTimestamp
				}
				endTime := now.Add(-delay)

				if startTime.After(endTime) {
					log.Warn("too soon to retrieve messages for topic", zap.String("topic", topic))
					continue
				}

				// Determine if the messages exist across all nodes
				for _, node := range options.StoreNodes {
					// TODO: make async
					storeNodeFailure := false
					var result *store.Result

				retry1:
					for i := 0; i < maxAttempts; i++ {
						result, err = wakuNode.Store().Query(ctx, store.FilterCriteria{
							ContentFilter: protocol.NewContentFilter(topic),
							TimeStart:     proto.Int64(startTime.UnixNano()),
							TimeEnd:       proto.Int64(endTime.UnixNano()),
						}, store.WithPeerAddr(node), store.IncludeData(false))
						if err != nil {
							logger.Error("could not query storenode", zap.Stringer("storenode", node), zap.Error(err))
							storeNodeFailure = true
							time.Sleep(2 * time.Second)
						} else {
							storeNodeFailure = false
							break retry1
						}
					}

					if storeNodeFailure {
						// TODO: Notify that storenode was not available from X to Y time
						logger.Error("storenode not available", zap.Stringer("storenode", node), zap.Time("startTime", startTime), zap.Time("endTime", endTime))
					} else {
						for {
							storeNodeFailure = false
							var hasNext bool
						retry2:
							for i := 0; i < maxAttempts; i++ {
								hasNext, err = result.Next(ctx)
								if err != nil {
									logger.Error("could not query storenode", zap.Stringer("storenode", node), zap.Error(err))
									storeNodeFailure = true
									time.Sleep(2 * time.Second)
								} else {
									break retry2
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
							} else {
								if !hasNext { // No more messages available
									break
								}

								for _, mkv := range result.Messages() {
									hash := mkv.WakuMessageHash()
									_, ok := msgMap[hash]
									if !ok {
										msgMap[hash] = make(map[string]MessageExistence)
									}
									msgMap[hash][node.String()] = Exists
									msgTopic[hash] = mkv.PubsubTopic
								}
							}
						}
					}
				}

				// Update db with last sync time
				dbStore.UpdateTopicSyncState(trx, options.ClusterID, topic, endTime)
			}

			// Verify for each storenode which messages are not available, and query for their existence using message hash

			// Node -> msgHash
			msgsToVerify := make(map[string][]pb.MessageHash)
			for msgHash, nodes := range msgMap {
				for node, existence := range nodes {
					if existence != Exists {
						msgsToVerify[node] = append(msgsToVerify[node], msgHash)
					}
				}
			}

			// TODO: async
			for node, messageHashes := range msgsToVerify {
				nodeMultiaddr, err := multiaddr.NewMultiaddr(node)
				if err != nil {
					return err
				}

				storeNodeFailure := false
				var result *store.Result
			retry3:
				for i := 0; i < maxAttempts; i++ {
					result, err = wakuNode.Store().QueryByHash(ctx, messageHashes, store.IncludeData(false), store.WithPeerAddr(nodeMultiaddr))
					if err != nil {
						logger.Error("could not query storenode", zap.Stringer("storenode", nodeMultiaddr), zap.Error(err))
						storeNodeFailure = true
						time.Sleep(2 * time.Second)
					} else {
						break retry3
					}
				}

				if storeNodeFailure {
					// TODO: Notify that storenode was not available from X to Y time
					logger.Error("storenode not available",
						zap.String("storenode", node),
						zap.Any("hashes", messageHashes))
				} else {
					for {
						storeNodeFailure = false
						var hasNext bool
					retry4:
						for i := 0; i < maxAttempts; i++ {
							hasNext, err = result.Next(ctx)
							if err != nil {
								logger.Error("could not query storenode", zap.String("storenode", node), zap.Error(err))
								storeNodeFailure = true
								time.Sleep(2 * time.Second)
							} else {
								break retry4
							}
						}

						if storeNodeFailure {
							// TODO: Notify that storenode was not available from X to Y time
							logger.Error("storenode not available",
								zap.String("storenode", node),
								zap.Any("hashes", messageHashes),
								zap.String("cursor", hexutil.Encode(result.Cursor())))
						} else {
							if !hasNext { // No more messages available
								break
							}

							for _, mkv := range result.Messages() {
								hash := mkv.WakuMessageHash()
								_, ok := msgMap[hash]
								if !ok {
									msgMap[hash] = make(map[string]MessageExistence)
								}
								msgMap[hash][node] = Exists
							}
						}
					}
				}
			}

			// TODO: if a message is not available, store in DB in which store nodes it wasnt available and from which time to which time
			for msgHash, nodes := range msgMap {
				var missingIn []string
				for node, existence := range nodes {
					if existence != Exists {
						missingIn = append(missingIn, node)
					}
				}

				err := dbStore.RecordMissingMessage(trx, msgHash, msgTopic[msgHash], missingIn)
				if err != nil {
					return err
				}
			}

			err = trx.Commit()
			// TODO: revert?
			if err != nil {
				return err
			}

		}
	}
}
