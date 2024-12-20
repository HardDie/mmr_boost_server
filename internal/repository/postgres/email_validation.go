//nolint:dupl
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/dimonrus/gosql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

type EmailValidation struct {
	db *db.DB
}

func NewEmailValidation(db *db.DB) *EmailValidation {
	return &EmailValidation{
		db: db,
	}
}

func (r *EmailValidation) CreateOrUpdate(ctx context.Context, userID int32, code string, expiredAt time.Time) (*entity.EmailValidation, error) {
	tx := getTxOrConn(ctx, r.db)

	ent := &entity.EmailValidation{
		UserID:    userID,
		Code:      strings.ToLower(code),
		ExpiredAt: expiredAt,
	}

	q := gosql.NewInsert().Into("email_validations")
	q.Columns().Add("user_id", "code", "expired_at")
	q.Columns().Arg(ent.UserID, ent.Code, ent.ExpiredAt)
	q.Conflict().Object("user_id").Action("UPDATE").Set().
		Add("code = EXCLUDED.code", "expired_at = EXCLUDED.expired_at", "created_at = now()")
	q.Returning().Add("id", "created_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&ent.ID, &ent.CreatedAt)
	if err != nil {
		logger.Error.Println("CreateOrUpdate:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return ent, nil
}
func (r *EmailValidation) GetByCode(ctx context.Context, code string) (*entity.EmailValidation, error) {
	tx := getTxOrConn(ctx, r.db)

	ent := &entity.EmailValidation{
		Code: strings.ToLower(code),
	}

	q := gosql.NewSelect().From("email_validations")
	q.Columns().Add("id", "user_id", "expired_at", "created_at")
	q.Where().AddExpression("code = ?", ent.Code)
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&ent.ID, &ent.UserID, &ent.ExpiredAt, &ent.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("GetByCode:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return ent, nil
}
func (r *EmailValidation) DeleteByID(ctx context.Context, id int32) error {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewDelete().From("email_validations")
	q.Where().AddExpression("id = ?", id)
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetGetArguments()...)

	err := row.Scan(&id)
	if err != nil {
		logger.Error.Println("DeleteByID:", err.Error())
		return status.Error(codes.Internal, "internal")
	}
	return nil
}
