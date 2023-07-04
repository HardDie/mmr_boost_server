package postgres_test

import (
	"testing"

	"github.com/lib/pq"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/migration"
)

var (
	defaultDB  = "db"
	psqlConfig = config.Postgres{
		Host:     "localhost",
		Port:     5432,
		User:     "mmr_boost",
		Password: "mmr_boost",
		Database: "test_db",
	}
)

func cleanDB(t *testing.T, dbConn *db.DB) {
	t.Helper()

	_, err := dbConn.DB.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("users DELETE err = %v; want nil", err)
	}
}

func setup(t *testing.T) *db.DB {
	t.Helper()

	createTestDB(t)

	// Connect to test_db
	dbConn, err := db.Get(psqlConfig)
	if err != nil {
		t.Fatalf("db.Get() err = %v; want nil", err)
	}
	t.Cleanup(func() {
		dbConn.DB.Close()
		dropTestDB(t)
	})

	err = migration.NewMigrate(dbConn).Up()
	if err != nil {
		t.Fatalf("Up() err = %v; want nil", err)
	}
	return dbConn
}

func createTestDB(t *testing.T) {
	cfg := psqlConfig
	cfg.Database = defaultDB
	dbConn, err := db.Get(cfg)
	if err != nil {
		t.Fatalf("db.Get() err = %v; want nil", err)
	}
	defer dbConn.DB.Close()
	_, err = dbConn.DB.Exec("CREATE DATABASE test_db")
	if err != nil {
		e, ok := err.(*pq.Error)
		if !ok {
			t.Fatalf("db CREATE err = %T; want *pq.Error", err)
		}
		if string(e.Code) != "42P04" {
			t.Fatalf("db CREATE err = %v; want nil", err)
		}
		// db already exists, skip error
	}
}
func dropTestDB(t *testing.T) {
	// Drop test_db
	cfg := psqlConfig
	cfg.Database = defaultDB
	dbConn, err := db.Get(cfg)
	if err != nil {
		t.Fatalf("db.Get() err = %v; want nil", err)
	}
	defer dbConn.DB.Close()
	_, err = dbConn.DB.Exec("DROP DATABASE test_db")
	if err != nil {
		t.Fatalf("db DROP err = %v; want nil", err)
	}
}
