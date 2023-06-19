package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dimonrus/gosql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type User struct {
	db *db.DB
}

func NewUser(db *db.DB) *User {
	return &User{
		db: db,
	}
}

func (r *User) GetByID(ctx context.Context, id int32) (*entity.User, error) { //nolint:dupl
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		ID: id,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("email", "username", "role_id", "steam_id", "is_activated", "created_at", "updated_at", "deleted_at")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.Email, &u.Username, &u.RoleID, &u.SteamID, &u.IsActivated, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("GetByID:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return u, nil
}
func (r *User) GetByName(ctx context.Context, name string) (*entity.User, error) { //nolint:dupl
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		Username: name,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("id", "email", "role_id", "steam_id", "is_activated", "created_at", "updated_at", "deleted_at")
	q.Where().AddExpression("username = ?", name)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.ID, &u.Email, &u.RoleID, &u.SteamID, &u.IsActivated, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("GetByName:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return u, nil
}
func (r *User) GetByNameOrEmail(ctx context.Context, name string, email string) (*entity.User, error) {
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		Username: name,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("id", "email", "role_id", "steam_id", "is_activated", "created_at", "updated_at", "deleted_at")
	q.Where().Merge(gosql.ConditionOperatorAnd,
		gosql.NewSqlCondition(gosql.ConditionOperatorOr).
			AddExpression("username = ?", name).
			AddExpression("email = ?", email),
	)
	q.Where().AddExpression("deleted_at IS NULL")
	q.SetPagination(1, 0)

	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.ID, &u.Email, &u.RoleID, &u.SteamID, &u.IsActivated, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("GetByNameOrEmail:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return u, nil
}
func (r *User) Create(ctx context.Context, email, name string) (*entity.User, error) {
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		Email:       email,
		Username:    name,
		RoleID:      int32(pb.UserRoleID_user),
		IsActivated: false,
	}

	q := gosql.NewInsert().Into("users")
	q.Columns().Add("email", "username", "role_id", "is_activated")
	q.Columns().Arg(u.Email, u.Username, u.RoleID, u.IsActivated)
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		logger.Error.Println("Create:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return u, nil
}
func (r *User) ActivateRecord(ctx context.Context, userID int32) (*entity.User, error) {
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		ID:          userID,
		IsActivated: true,
	}

	q := gosql.NewUpdate().Table("users")
	q.Set().Append("is_activated = true")
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("email", "username", "role_id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.Email, &u.Username, &u.RoleID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error.Println("ActivateRecord:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return u, nil
}
func (r *User) UpdateSteamID(ctx context.Context, userID int32, steamID string) (*entity.User, error) {
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		ID:      userID,
		SteamID: &steamID,
	}

	q := gosql.NewUpdate().Table("users")
	q.Set().Append("steam_id = ?", u.SteamID)
	q.Set().Append("updated_at = now()")
	q.Where().AddExpression("id = ?", u.ID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("email", "username", "role_id", "is_activated", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.Email, &u.Username, &u.RoleID, &u.IsActivated, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		logger.Error.Println("UpdateSteamID:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	return u, nil
}
