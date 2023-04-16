package postgres

import (
	"context"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
)

type application struct {
	db *db.DB
}

func newApplication(db *db.DB) application {
	return application{
		db: db,
	}
}

func (r *application) ApplicationCreate(ctx context.Context, req *dto.ApplicationCreateRequest, userID int32) (*entity.ApplicationPublic, error) {
	tx := getTxOrConn(ctx, r.db)

	app := &entity.ApplicationPublic{
		UserID:     userID,
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
