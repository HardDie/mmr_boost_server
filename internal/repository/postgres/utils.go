package postgres

import (
	"context"

	"github.com/HardDie/godb/v2"

	"github.com/HardDie/mmr_boost_server/internal/db"
)

type txKey struct{}

func setCtxTx(ctx context.Context, tx *godb.SqlTx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func getCtxTx(ctx context.Context) *godb.SqlTx {
	value := ctx.Value(txKey{})
	if value == nil {
		return nil
	}
	tx, ok := value.(*godb.SqlTx)
	if !ok {
		return nil
	}
	return tx
}

func getTxOrConn(ctx context.Context, db *db.DB) godb.Queryer {
	if tx := getCtxTx(ctx); tx != nil {
		return tx
	}
	return db.DB
}
