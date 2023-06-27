package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dimonrus/gosql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type Application struct {
	db *db.DB
}

func NewApplication(db *db.DB) *Application {
	return &Application{
		db: db,
	}
}

func (r *Application) Create(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error) {
	tx := getTxOrConn(ctx, r.db)

	app := &entity.ApplicationPublic{
		UserID:     req.UserID,
		StatusID:   int32(pb.ApplicationStatusID_created),
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
		logger.Error.Println("Create:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return app, nil
}

func (r *Application) List(ctx context.Context, req *dto.ApplicationListRequest) ([]*entity.ApplicationPublic, error) {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewSelect().From("applications")
	q.Columns().Add("id", "user_id", "status_id", "type_id", "current_mmr", "target_mmr", "tg_contact",
		"created_at", "updated_at")
	if req.UserID != nil {
		q.Where().AddExpression("user_id = ?", req.UserID)
	}
	if req.StatusID != nil {
		q.Where().AddExpression("status_id = ?", req.StatusID)
	}
	q.Where().AddExpression("deleted_at IS NULL")
	q.AddOrder("id DESC")
	rows, err := tx.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		logger.Error.Println("error select applications from DB:", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	defer rows.Close()

	var res []*entity.ApplicationPublic
	for rows.Next() {
		app := &entity.ApplicationPublic{}
		err = rows.Scan(&app.ID, &app.UserID, &app.StatusID, &app.TypeID, &app.CurrentMMR, &app.TargetMMR, &app.TgContact,
			&app.CreatedAt, &app.UpdatedAt)
		if err != nil {
			logger.Error.Println("error scan applications row from DB:", err.Error())
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

func (r *Application) Item(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPublic, error) {
	tx := getTxOrConn(ctx, r.db)

	app := &entity.ApplicationPublic{
		ID: req.ApplicationID,
	}

	q := gosql.NewSelect().From("applications")
	q.Columns().Add("user_id", "status_id", "type_id", "current_mmr", "target_mmr", "tg_contact",
		"created_at", "updated_at")
	q.Where().AddExpression("id = ?", req.ApplicationID)
	if req.UserID != nil {
		q.Where().AddExpression("user_id = ?", req.UserID)
	}
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&app.UserID, &app.StatusID, &app.TypeID, &app.CurrentMMR, &app.TargetMMR, &app.TgContact,
		&app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("Item:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return app, nil
}
func (r *Application) PrivateItem(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPrivate, error) {
	tx := getTxOrConn(ctx, r.db)

	app := &entity.ApplicationPrivate{
		ID: req.ApplicationID,
	}

	q := gosql.NewSelect().From("applications")
	q.Columns().Add("steam_login", "steam_password", "created_at", "updated_at")
	q.Where().AddExpression("id = ?", req.ApplicationID)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&app.SteamLogin, &app.SteamPassword, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("PrivateItem:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return app, nil
}
