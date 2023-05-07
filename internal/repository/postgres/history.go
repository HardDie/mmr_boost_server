package postgres

import (
	"context"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/mmr_boost_server/internal/db"
)

type history struct {
	db *db.DB
}

func NewHistory(db *db.DB) *history {
	return &history{
		db: db,
	}
}

func (r *history) NewEvent(ctx context.Context, userID int32, message string) error {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewInsert().Into("history")
	q.Columns().Add("user_id", "message")
	q.Columns().Arg(userID, message)
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	var id int32
	err := row.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
func (r *history) NewEventWithAffected(ctx context.Context, userID, affectedUserID int32, message string) error {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewInsert().Into("history")
	q.Columns().Add("user_id", "affect_user_id", "message")
	q.Columns().Arg(userID, affectedUserID, message)
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	var id int32
	err := row.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
