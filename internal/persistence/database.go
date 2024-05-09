package persistence

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/timesource"
	"go.uber.org/zap"
)

// WALMode for sqlite.
const WALMode = "wal"

// DBStore is a MessageProvider that has a *sql.DB connection
type DBStore struct {
	db          *sql.DB
	migrationFn func(db *sql.DB, logger *zap.Logger) error

	timesource timesource.Timesource
	log        *zap.Logger

	enableMigrations bool

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// DBOption is an optional setting that can be used to configure the DBStore
type DBOption func(*DBStore) error

// WithDB is a DBOption that lets you use any custom *sql.DB with a DBStore.
func WithDB(db *sql.DB) DBOption {
	return func(d *DBStore) error {
		d.db = db
		return nil
	}
}

// ConnectionPoolOptions is the options to be used for DB connection pooling
type ConnectionPoolOptions struct {
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

// WithDriver is a DBOption that will open a *sql.DB connection
func WithDriver(driverName string, datasourceName string, connectionPoolOptions ...ConnectionPoolOptions) DBOption {
	return func(d *DBStore) error {
		db, err := sql.Open(driverName, datasourceName)
		if err != nil {
			return err
		}

		if len(connectionPoolOptions) != 0 {
			db.SetConnMaxIdleTime(connectionPoolOptions[0].ConnectionMaxIdleTime)
			db.SetConnMaxLifetime(connectionPoolOptions[0].ConnectionMaxLifetime)
			db.SetMaxIdleConns(connectionPoolOptions[0].MaxIdleConnections)
			db.SetMaxOpenConns(connectionPoolOptions[0].MaxOpenConnections)
		}

		d.db = db
		return nil
	}
}

type MigrationFn func(db *sql.DB, logger *zap.Logger) error

// WithMigrations is a DBOption used to determine if migrations should
// be executed, and what driver to use
func WithMigrations(migrationFn MigrationFn) DBOption {
	return func(d *DBStore) error {
		d.enableMigrations = true
		d.migrationFn = migrationFn
		return nil
	}
}

// DefaultOptions returns the default DBoptions to be used.
func DefaultOptions() []DBOption {
	return []DBOption{}
}

// Creates a new DB store using the db specified via options.
// It will run migrations if enabled
// clean up records according to the retention policy used
func NewDBStore(log *zap.Logger, options ...DBOption) (*DBStore, error) {
	result := new(DBStore)
	result.log = log.Named("dbstore")

	optList := DefaultOptions()
	optList = append(optList, options...)

	for _, opt := range optList {
		err := opt(result)
		if err != nil {
			return nil, err
		}
	}

	if result.enableMigrations {
		err := result.migrationFn(result.db, log)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Start starts the store server functionality
func (d *DBStore) Start(ctx context.Context, timesource timesource.Timesource) error {
	ctx, cancel := context.WithCancel(ctx)

	d.cancel = cancel
	d.timesource = timesource

	err := d.cleanOlderRecords(ctx)
	if err != nil {
		return err
	}

	d.wg.Add(1)
	go d.checkForOlderRecords(ctx, 60*time.Second)

	return nil
}

func (d *DBStore) cleanOlderRecords(ctx context.Context) error {
	d.log.Info("Cleaning older records...")

	// TODO:
	deleteFrom := time.Now().Add(-14 * 24 * time.Hour)
	_, err := d.db.Exec("DELETE FROM missingMessages WHERE storedAt <", deleteFrom)
	if err != nil {
		return err
	}

	d.log.Info("Older records removed")

	return nil
}

func (d *DBStore) checkForOlderRecords(ctx context.Context, t time.Duration) {
	defer d.wg.Done()

	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := d.cleanOlderRecords(ctx)
			if err != nil {
				d.log.Error("cleaning older records", zap.Error(err))
			}
		}
	}
}

// Stop closes a DB connection
func (d *DBStore) Stop() {
	if d.cancel == nil {
		return
	}

	d.cancel()
	d.wg.Wait()
	d.db.Close()
}

func (d *DBStore) GetTrx(ctx context.Context) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, &sql.TxOptions{})
}

func (d *DBStore) GetTopicSyncStatus(ctx context.Context, clusterID uint, pubsubTopics []string) (map[string]*time.Time, error) {
	result := make(map[string]*time.Time)
	for _, topic := range pubsubTopics {
		result[topic] = &time.Time{}
	}

	sqlQuery := `SELECT pubsubTopic, lastSyncTimestamp FROM syncTopicStatus WHERE clusterId = ?`
	rows, err := d.db.QueryContext(ctx, sqlQuery, clusterID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pubsubTopic string
		var lastSyncTimestamp int64
		err := rows.Scan(&pubsubTopic, &lastSyncTimestamp)
		if err != nil {
			return nil, err
		}

		t := time.Unix(0, lastSyncTimestamp)
		result[pubsubTopic] = &t
	}
	defer rows.Close()

	return result, nil
}

func (d *DBStore) UpdateTopicSyncState(tx *sql.Tx, clusterID uint, topic string, lastSyncTimestamp time.Time) error {
	stmt, err := tx.Prepare("INSERT INTO syncTopicStatus(clusterId, pubsubTopic, lastSyncTimestamp) VALUES (?, ?, ?) ON CONFLICT(clusterId, pubsubTopic) DO UPDATE SET lastSyncTimestamp = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(clusterID, topic, lastSyncTimestamp.UnixNano())
	if err != nil {
		return err
	}

	return stmt.Close()
}

func (d *DBStore) RecordMessage(tx *sql.Tx, msgHash pb.MessageHash, clusterID uint, topic string, timestamp uint64, storenodes []string, status string) error {
	stmt, err := tx.Prepare("INSERT INTO missingMessages(clusterId, pubsubTopic, messageHash, msgTimestamp, storenode, status, storedAt) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, s := range storenodes {
		_, err := stmt.Exec(clusterID, topic, msgHash.String(), timestamp, s, status, now)
		if err != nil {
			return err
		}
	}

	return stmt.Close()
}
