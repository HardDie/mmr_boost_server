package postgres

import (
	"context"

	"github.com/dimonrus/gosql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

type History struct {
	db *db.DB
}

func NewHistory(db *db.DB) *History {
	return &History{
		db: db,
	}
}

func (r *History) NewEvent(ctx context.Context, userID int32, message string) error {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewInsert().Into("history")
	q.Columns().Add("user_id", "message")
	q.Columns().Arg(userID, message)
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	var id int32
	err := row.Scan(&id)
	if err != nil {
		logger.Error.Println("NewEvent:", err.Error())
		return status.Error(codes.Internal, "internal")
	}

	return nil
}
func (r *History) NewEventWithAffected(ctx context.Context, userID, affectedUserID int32, message string) error {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewInsert().Into("history")
	q.Columns().Add("user_id", "affect_user_id", "message")
	q.Columns().Arg(userID, affectedUserID, message)
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	var id int32
	err := row.Scan(&id)
	if err != nil {
		logger.Error.Println("NewEventWithAffected:", err.Error())
		return status.Error(codes.Internal, "internal")
	}

	return nil
}
