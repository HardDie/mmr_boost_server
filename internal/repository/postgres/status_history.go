package postgres

import (
	"context"

	"github.com/dimonrus/gosql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

type StatusHistory struct {
	db *db.DB
}

func NewStatusHistory(db *db.DB) *StatusHistory {
	return &StatusHistory{
		db: db,
	}
}

func (r *StatusHistory) NewEvent(ctx context.Context, req *dto.StatusHistoryNewEventRequest) error {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewInsert().Into("status_history")
	q.Columns().Add("user_id", "application_id", "new_status_id")
	q.Columns().Arg(req.UserID, req.ApplicationID, req.NewStatusID)
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
func (r *StatusHistory) List(ctx context.Context, applicationID int32) ([]*entity.StatusHistory, error) {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewSelect().From("status_history")
	q.Columns().Add("id", "user_id", "application_id", "new_status_id", "created_at")
	q.Where().AddExpression("application_id = ?", applicationID)
	q.AddOrder("id")
	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		logger.Error.Println("error select status_history from DB:", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	defer rows.Close()

	var res []*entity.StatusHistory
	for rows.Next() {
		app := &entity.StatusHistory{}
		err = rows.Scan(&app.ID, &app.UserID, &app.ApplicationID, &app.NewStatusID, &app.CreatedAt)
		if err != nil {
			logger.Error.Println("error scan status_history row from DB:", err.Error())
			return nil, status.Error(codes.Internal, "internal error")
		}
		res = append(res, app)
	}
	if err = rows.Err(); err != nil {
		logger.Error.Println("rows error:", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return res, nil
}
