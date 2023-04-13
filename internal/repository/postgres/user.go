package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/entity"
)

type user struct {
	db *db.DB
}

func newUser(db *db.DB) user {
	return user{
		db: db,
	}
}

func (r *user) UserGetByID(ctx context.Context, id int32) (*entity.User, error) {
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		ID: id,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("email", "username", "role_id", "is_activated", "created_at", "updated_at", "deleted_at")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.Email, &u.Username, &u.RoleID, &u.IsActivated, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil

}
func (r *user) UserGetByName(ctx context.Context, name string) (*entity.User, error) {
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		Username: name,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("id", "email", "role_id", "is_activated", "created_at", "updated_at", "deleted_at")
	q.Where().AddExpression("username = ?", name)
	q.Where().AddExpression("deleted_at IS NULL")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.ID, &u.Email, &u.RoleID, &u.IsActivated, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}
func (r *user) UserCreate(ctx context.Context, email, name string) (*entity.User, error) {
	tx := getTxOrConn(ctx, r.db)

	u := &entity.User{
		Email:       email,
		Username:    name,
		RoleID:      4,
		IsActivated: false,
	}

	q := gosql.NewInsert().Into("users")
	q.Columns().Add("email", "username", "role_id", "is_activated")
	q.Columns().Arg(u.Email, u.Username, u.RoleID, u.IsActivated)
	q.Returning().Add("id", "created_at", "updated_at")
	row := tx.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
