package postgres

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/jackc/pgx/v5/stdlib" // Blank import to register the postgres driver
	"github.com/waku-org/storenode-messages/internal/persistence/migrate"
	"github.com/waku-org/storenode-messages/internal/persistence/postgres/migrations"
	"go.uber.org/zap"
)

// NewDB connects to postgres DB in the specified path
func NewDB(dburl string, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dburl)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrationDriver(db *sql.DB) (database.Driver, error) {
	return pgx.WithInstance(db, &pgx.Config{
		MigrationsTable: pgx.DefaultMigrationsTable,
	})
}

// Migrations is the function used for DB migration with postgres driver
func Migrations(db *sql.DB, logger *zap.Logger) error {
	migrationDriver, err := migrationDriver(db)
	if err != nil {
		return err
	}
	return migrate.Migrate(db, migrationDriver, migrations.AssetNames(), migrations.Asset)
}
