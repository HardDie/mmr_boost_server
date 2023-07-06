package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dimonrus/gosql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

type AccessToken struct {
	db      *db.DB
	timeNow func() time.Time
}

func NewAccessToken(db *db.DB) *AccessToken {
	return &AccessToken{
		db:      db,
		timeNow: time.Now,
	}
}

func (r *AccessToken) CreateOrUpdate(ctx context.Context, userID int32, tokenHash string, expiredAt time.Time) (*entity.AccessToken, error) {
	tx := getTxOrConn(ctx, r.db)

	token := &entity.AccessToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiredAt: expiredAt,
		CreatedAt: r.timeNow(),
		UpdatedAt: r.timeNow(),
	}

	q := gosql.NewInsert().Into("access_tokens")
	q.Columns().Add("user_id", "token_hash", "created_at", "updated_at", "expired_at")
	q.Columns().Arg(token.UserID, token.TokenHash, token.CreatedAt, token.UpdatedAt, token.ExpiredAt)
	q.Conflict().Object("user_id").Action("UPDATE").Set().
		Add("token_hash = EXCLUDED.token_hash", "expired_at = EXCLUDED.expired_at", "updated_at = EXCLUDED.updated_at", "deleted_at = NULL")
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&token.ID, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		logger.Error.Println("CreateOrUpdate:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return token, nil
}
func (r *AccessToken) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.AccessToken, error) {
	tx := getTxOrConn(ctx, r.db)

	token := &entity.AccessToken{
		TokenHash: tokenHash,
	}

	q := gosql.NewSelect().From("access_tokens")
	q.Columns().Add("id", "user_id", "expired_at", "created_at", "updated_at")
	q.Where().AddExpression("token_hash = ?", token.TokenHash)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&token.ID, &token.UserID, &token.ExpiredAt, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("GetByUserID:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return token, nil
}
func (r *AccessToken) DeleteByID(ctx context.Context, id int32) error {
	tx := getTxOrConn(ctx, r.db)

	q := gosql.NewUpdate().Table("access_tokens")
	q.Set().Append("deleted_at = ?", r.timeNow())
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&id)
	if err != nil {
		logger.Error.Println("DeleteByID:", err.Error())
		return status.Error(codes.Internal, "internal")
	}
	return nil
}
