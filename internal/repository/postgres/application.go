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

type application struct {
	db *db.DB
}

func newApplication(db *db.DB) application {
	return application{
		db: db,
	}
}

func (r *application) ApplicationCreate(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error) {
	tx := getTxOrConn(ctx, r.db)

	app := &entity.ApplicationPublic{
		UserID:     req.UserID,
		StatusID:   1,
		TypeID:     req.TypeID,
		CurrentMMR: req.CurrentMMR,
		TargetMMR:  req.TargetMMR,
		TgContact:  req.TgContact,
	}

	q := gosql.NewInsert().Into("applications")
	q.Columns().Add("user_id", "status_id", "type_id", "current_mmr", "target_mmr", "tg_contact")
	q.Columns().Arg(app.UserID, app.StatusID, app.TypeID, app.CurrentMMR, app.TargetMMR, app.TgContact)
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (r *application) ApplicationUserList(ctx context.Context, req *dto.ApplicationUserListRequest) ([]*entity.ApplicationPublic, error) {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewSelect().From("applications")
	q.Columns().Add("id", "status_id", "type_id", "current_mmr", "target_mmr", "tg_contact", "created_at", "updated_at")
	q.Where().AddExpression("user_id = ?", req.UserID)
	if req.StatusID != nil {
		q.Where().AddExpression("status_id = ?", req.StatusID)
	}
	q.Where().AddExpression("deleted_at IS NULL")
	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		logger.Error.Println("error select applications from DB:", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}

	var res []*entity.ApplicationPublic
	for rows.Next() {
		app := &entity.ApplicationPublic{
			UserID: req.UserID,
		}
		err = rows.Scan(&app.ID, &app.StatusID, &app.TypeID, &app.CurrentMMR, &app.TargetMMR, &app.TgContact,
			&app.CreatedAt, &app.UpdatedAt)
		if err != nil {
			logger.Error.Println("error scan applications row from DB:", err.Error())
			return nil, status.Error(codes.Internal, "internal error")
		}
		res = append(res, app)
	}

	return res, nil
}
