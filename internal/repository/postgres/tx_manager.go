package postgres

import (
	"context"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

type TxManager interface {
	ReadWriteTx(context.Context, func(ctx context.Context) error) error
}

type txManager struct {
	db *db.DB
}

func newTxManager(db *db.DB) *txManager {
	return &txManager{
		db: db,
	}
}

func (t *txManager) ReadWriteTx(ctx context.Context, call func(ctx context.Context) error) error {
	// Start transaction
	tx, err := t.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Close transaction at exit
	defer func() {
		txErr := t.db.EndTx(tx, err)
		if txErr != nil {
			logger.Error.Printf("Error close transaction: %v\n", txErr)
		}
	}()

	// Put transaction into context
	ctx = setCtxTx(ctx, tx)

	// Proceed
	err = call(ctx)
	if err != nil {
		return err
	}

	return nil
}
