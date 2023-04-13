package migration

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/migrations"
)

const (
	MigrationTable = "migrations"
)

type Migrate struct {
	db *sql.DB
}

func NewMigrate(db *db.DB) *Migrate {
	goose.SetBaseFS(migrations.Migrations)
	goose.SetTableName(MigrationTable)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error.Fatal(err)
	}

	return &Migrate{db: db.DB.DB}
}

func (m *Migrate) Up() error {
	err := goose.Up(m.db, ".")
	if err != nil {
		return fmt.Errorf("migrations failed: %w", err)
	}
	return nil
}
