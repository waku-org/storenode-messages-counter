package sqlite

import (
	"database/sql"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3" // Blank import to register the sqlite3 driver
	"github.com/waku-org/storenode-messages/internal/persistence/migrate"
	"github.com/waku-org/storenode-messages/internal/persistence/sqlite/migrations"
	"go.uber.org/zap"
)

func addSqliteURLDefaults(dburl string) string {
	if !strings.Contains(dburl, "?") {
		dburl += "?"
	}

	if !strings.Contains(dburl, "_journal=") {
		dburl += "&_journal=WAL"
	}

	if !strings.Contains(dburl, "_timeout=") {
		dburl += "&_timeout=5000"
	}

	return dburl
}

// NewDB creates a sqlite3 DB in the specified path
func NewDB(dburl string, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", addSqliteURLDefaults(dburl))
	if err != nil {
		return nil, err
	}

	// Disable concurrent access as not supported by the driver
	db.SetMaxOpenConns(1)

	return db, nil
}

func migrationDriver(db *sql.DB) (database.Driver, error) {
	return sqlite3.WithInstance(db, &sqlite3.Config{
		MigrationsTable: "message_counter_" + sqlite3.DefaultMigrationsTable,
	})
}

// Migrations is the function used for DB migration with sqlite driver
func Migrations(db *sql.DB, logger *zap.Logger) error {
	migrationDriver, err := migrationDriver(db)
	if err != nil {
		return err
	}

	return migrate.Migrate(db, migrationDriver, migrations.AssetNames(), migrations.Asset)
}
