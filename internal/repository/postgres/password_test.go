//go:build int
// +build int

package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/utils"
)

func comparePassword(t *testing.T, a, b *entity.Password) {
	t.Helper()

	if a == nil && b == nil {
		return
	}
	if reflect.DeepEqual(a, b) {
		return
	}
	if a.ID != b.ID {
		t.Errorf("ID = %d; want %d", a.ID, b.ID)
	}
	if a.UserID != b.UserID {
		t.Errorf("UserID = %d; want %d", a.UserID, b.UserID)
	}
	if a.PasswordHash != b.PasswordHash {
		t.Errorf("PasswordHash = %q; want %q", a.PasswordHash, b.PasswordHash)
	}
	if a.FailedAttempts != b.FailedAttempts {
		t.Errorf("FailedAttempts = %d; want %d", a.FailedAttempts, b.FailedAttempts)
	}
	if !a.CreatedAt.Equal(b.CreatedAt) {
		t.Errorf("CreatedAt = %v; want %v", a.CreatedAt, b.CreatedAt)
	}
	if !a.UpdatedAt.Equal(b.UpdatedAt) {
		t.Errorf("UpdatedAt = %v; want %v", a.UpdatedAt, b.UpdatedAt)
	}
	if !utils.CompareTime(a.DeletedAt, b.DeletedAt) {
		t.Errorf("DeletedAt = %v; want %v", a.DeletedAt, b.DeletedAt)
	}
	if !utils.CompareTime(a.BlockedAt, b.BlockedAt) {
		t.Errorf("BlockedAt = %v; want %v", a.BlockedAt, b.BlockedAt)
	}
}

func TestRepositoryPassword(t *testing.T) {
	dbConn := setup(t)

	t.Run("GetByUserID", testRepositoryPassword_GetByUserID(dbConn))
	t.Run("Create", testRepositoryPassword_Create(dbConn))
	t.Run("Update", testRepositoryPassword_Update(dbConn))
	t.Run("IncreaseFailedAttempts", testRepositoryPassword_IncreaseFailedAttempts(dbConn))
	t.Run("ResetFailedAttempts", testRepositoryPassword_ResetFailedAttempts(dbConn))
}

func testRepositoryPassword_GetByUserID(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(userID int32, db *db.DB) *entity.Password
		}{
			"manually created password": {
				setup: func(userID int32, db *db.DB) *entity.Password {
					pwd := &entity.Password{
						UserID:         userID,
						PasswordHash:   "test",
						FailedAttempts: 1,
						CreatedAt:      now,
						UpdatedAt:      now,
					}

					q := gosql.NewInsert().Into("passwords")
					q.Columns().Add("user_id", "password_hash", "failed_attempts", "created_at", "updated_at")
					q.Columns().Arg(pwd.UserID, pwd.PasswordHash, pwd.FailedAttempts, pwd.CreatedAt, pwd.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&pwd.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}

					return pwd
				},
			},
			"not exist password": {
				setup: func(userID int32, db *db.DB) *entity.Password {
					return nil
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				passwordRepository := postgres.NewPassword(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				want := tc.setup(user.ID, dbConn)

				got, err := passwordRepository.GetByUserID(ctx, user.ID)
				if err != nil {
					t.Fatalf("GetByUserID() err = %v; want nil", err)
				}
				comparePassword(t, got, want)
			})
		}
	}
}
func testRepositoryPassword_Create(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			want *entity.Password
		}{
			"valid": {
				want: &entity.Password{
					PasswordHash: "test",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				passwordRepository := postgres.NewPassword(dbConn)
				passwordRepository.SetTimeNow(func() time.Time {
					return now
				})

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}
				tc.want.UserID = user.ID

				got, err := passwordRepository.Create(ctx, user.ID, tc.want.PasswordHash)
				if err != nil {
					t.Fatalf("password.Create() err = %v; want nil", err)
				}
				tc.want.ID = got.ID
				comparePassword(t, got, tc.want)

				got, err = passwordRepository.GetByUserID(ctx, user.ID)
				if err != nil {
					t.Fatalf("GetByUserID() err = %v; want nil", err)
				}
				comparePassword(t, got, tc.want)
			})
		}
	}
}
func testRepositoryPassword_Update(dbConn *db.DB) func(t *testing.T) {
	ctx := context.Background()
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password
		}{
			"valid": {
				setup: func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password {
					pwd := &entity.Password{
						UserID:       userID,
						PasswordHash: "new_test",
						CreatedAt:    now,
						UpdatedAt:    now.Add(time.Minute),
					}

					passwordRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := passwordRepository.Create(ctx, userID, "test")
					if err != nil {
						t.Fatalf("password.Create() err = %v; want nil", err)
					}
					pwd.ID = resp.ID

					return pwd
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				passwordRepository := postgres.NewPassword(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				want := tc.setup(t, passwordRepository, user.ID)

				passwordRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})
				got, err := passwordRepository.Update(ctx, want.ID, want.PasswordHash)
				if err != nil {
					t.Fatalf("Update() err = %v; want nil", err)
				}
				comparePassword(t, got, want)

				got, err = passwordRepository.GetByUserID(ctx, user.ID)
				if err != nil {
					t.Fatalf("GetByUserID() err = %v; want nil", err)
				}
				comparePassword(t, got, want)
			})
		}
	}
}
func testRepositoryPassword_IncreaseFailedAttempts(dbConn *db.DB) func(t *testing.T) {
	ctx := context.Background()
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password
		}{
			"valid": {
				setup: func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password {
					pwd := &entity.Password{
						UserID:         userID,
						PasswordHash:   "test",
						FailedAttempts: 1,
						CreatedAt:      now,
						UpdatedAt:      now.Add(time.Minute),
					}

					passwordRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := passwordRepository.Create(ctx, pwd.UserID, pwd.PasswordHash)
					if err != nil {
						t.Fatalf("password.Create() err = %v; want nil", err)
					}
					pwd.ID = resp.ID

					return pwd
				},
			},
			"with value": {
				setup: func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password {
					pwd := &entity.Password{
						UserID:         userID,
						PasswordHash:   "test",
						FailedAttempts: 2,
						CreatedAt:      now,
						UpdatedAt:      now.Add(time.Minute),
					}

					passwordRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := passwordRepository.Create(ctx, pwd.UserID, pwd.PasswordHash)
					if err != nil {
						t.Fatalf("password.Create() err = %v; want nil", err)
					}
					pwd.ID = resp.ID

					_, err = passwordRepository.IncreaseFailedAttempts(ctx, pwd.ID)
					if err != nil {
						t.Fatalf("IncreaseFailedAttempts() err = %v; want nil", err)
					}

					return pwd
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				passwordRepository := postgres.NewPassword(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				want := tc.setup(t, passwordRepository, user.ID)

				passwordRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})
				got, err := passwordRepository.IncreaseFailedAttempts(ctx, want.ID)
				if err != nil {
					t.Fatalf("IncreaseFailedAttempts() err = %v; want nil", err)
				}
				comparePassword(t, got, want)

				got, err = passwordRepository.GetByUserID(ctx, user.ID)
				if err != nil {
					t.Fatalf("GetByUserID() err = %v; want nil", err)
				}
				comparePassword(t, got, want)
			})
		}
	}
}
func testRepositoryPassword_ResetFailedAttempts(dbConn *db.DB) func(t *testing.T) {
	ctx := context.Background()
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password
		}{
			"valid": {
				setup: func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password {
					pwd := &entity.Password{
						UserID:         userID,
						PasswordHash:   "test",
						FailedAttempts: 0,
						CreatedAt:      now,
						UpdatedAt:      now.Add(time.Minute),
					}

					passwordRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := passwordRepository.Create(ctx, pwd.UserID, pwd.PasswordHash)
					if err != nil {
						t.Fatalf("password.Create() err = %v; want nil", err)
					}
					pwd.ID = resp.ID

					_, err = passwordRepository.IncreaseFailedAttempts(ctx, pwd.ID)
					if err != nil {
						t.Fatalf("IncreaseFailedAttempts() err = %v; want nil", err)
					}

					return pwd
				},
			},
			"empty": {
				setup: func(t *testing.T, passwordRepository *postgres.Password, userID int32) *entity.Password {
					pwd := &entity.Password{
						UserID:         userID,
						PasswordHash:   "test",
						FailedAttempts: 0,
						CreatedAt:      now,
						UpdatedAt:      now.Add(time.Minute),
					}

					passwordRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := passwordRepository.Create(ctx, pwd.UserID, pwd.PasswordHash)
					if err != nil {
						t.Fatalf("password.Create() err = %v; want nil", err)
					}
					pwd.ID = resp.ID

					return pwd
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				passwordRepository := postgres.NewPassword(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				want := tc.setup(t, passwordRepository, user.ID)

				passwordRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})
				got, err := passwordRepository.ResetFailedAttempts(ctx, want.ID)
				if err != nil {
					t.Fatalf("ResetFailedAttempts() err = %v; want nil", err)
				}
				comparePassword(t, got, want)

				got, err = passwordRepository.GetByUserID(ctx, user.ID)
				if err != nil {
					t.Fatalf("GetByUserID() err = %v; want nil", err)
				}
				comparePassword(t, got, want)
			})
		}
	}
}
