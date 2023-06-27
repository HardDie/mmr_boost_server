package db

import (
	"context"
	"fmt"
	"time"

	"github.com/HardDie/godb/v2"
	// import db driver.
	_ "github.com/lib/pq"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

const (
	MaxConnections         = 50
	ConnectionIdleLifetime = 15
)

type DB struct {
	DB *godb.DBO
}

func Get(cfg config.Postgres) (*DB, error) {
	conf := godb.PostgresConnectionConfig{
		ConnectionConfig: godb.ConnectionConfig{
			Host:                   cfg.Host,
			Port:                   cfg.Port,
			Name:                   cfg.Database,
			User:                   cfg.User,
			Password:               cfg.Password,
			MaxConnections:         MaxConnections,
			ConnectionIdleLifetime: ConnectionIdleLifetime,
		},
		SSLMode: "disable",
	}

	dbConfig := godb.DBO{
		Options: godb.Options{
			//Debug:  true,
			Logger: logger.Debug,
		},
		Connection: &conf,
	}

	var err error
	res := &DB{}

	// Reconnecting to the database in case of failure
	for i := 1; i < 8; i++ {
		res.DB, err = dbConfig.Init()
		if err == nil {
			break
		}
		logger.Error.Println("error open connection to db:", err.Error())
		time.Sleep(time.Duration(i*i) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("error open connection to db: %w", err)
	}

	return res, nil
}

func (db *DB) BeginTx(ctx context.Context) (*godb.SqlTx, error) {
	return db.DB.BeginContext(ctx)
}
func (db *DB) EndTx(tx *godb.SqlTx, err error) error {
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			logger.Error.Printf("error rollback tx: %v", err.Error())
			return err
		}
		return nil
	}

	err = tx.Commit()
	if err != nil {
		logger.Error.Printf("error commit tx: %v", err.Error())
		return err
	}
	return nil
}
