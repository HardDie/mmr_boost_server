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

func compareAccessToken(t *testing.T, a, b *entity.AccessToken) {
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
	if a.TokenHash != b.TokenHash {
		t.Errorf("TokenHash = %q; want %q", a.TokenHash, b.TokenHash)
	}
	if !a.ExpiredAt.Equal(b.ExpiredAt) {
		t.Errorf("ExpiredAt = %v; want %v", a.ExpiredAt, b.ExpiredAt)
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
}

func TestRepositoryAccessToken(t *testing.T) {
	dbConn := setup(t)

	t.Run("GetByUserID", testRepositoryAccessToken_GetByUserID(dbConn))
	t.Run("CreateOrUpdate", testRepositoryAccessToken_CreateOrUpdate(dbConn))
	t.Run("DeleteByID", testRepositoryAccessToken_DeleteByID(dbConn))
}

func testRepositoryAccessToken_GetByUserID(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(userID int32, db *db.DB) (string, *entity.AccessToken)
		}{
			"manually created access_token": {
				setup: func(userID int32, db *db.DB) (string, *entity.AccessToken) {
					token := &entity.AccessToken{
						UserID:    userID,
						TokenHash: "test",
						ExpiredAt: now.Add(time.Minute),
						CreatedAt: now,
						UpdatedAt: now,
					}

					q := gosql.NewInsert().Into("access_tokens")
					q.Columns().Add("user_id", "token_hash", "expired_at", "created_at", "updated_at")
					q.Columns().Arg(token.UserID, token.TokenHash, token.ExpiredAt, token.CreatedAt, token.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&token.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}

					return token.TokenHash, token
				},
			},
			"not exist access_token": {
				setup: func(userID int32, db *db.DB) (string, *entity.AccessToken) {
					return "invalid", nil
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				accessTokenRepository := postgres.NewAccessToken(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				token, want := tc.setup(user.ID, dbConn)

				got, err := accessTokenRepository.GetByTokenHash(ctx, token)
				if err != nil {
					t.Fatalf("GetByTokenHash() err = %v; want nil", err)
				}
				compareAccessToken(t, got, want)
			})
		}
	}
}
func testRepositoryAccessToken_CreateOrUpdate(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(dbConn *db.DB, accessTokenRepository *postgres.AccessToken, userID int32) *entity.AccessToken
		}{
			"valid": {
				setup: func(dbConn *db.DB, accessTokenRepository *postgres.AccessToken, userID int32) *entity.AccessToken {
					token := &entity.AccessToken{
						UserID:    userID,
						TokenHash: "test",
						ExpiredAt: now.Add(time.Minute),
						CreatedAt: now,
						UpdatedAt: now,
					}

					accessTokenRepository.SetTimeNow(func() time.Time {
						return now
					})
					return token
				},
			},
			"update exist": {
				setup: func(dbConn *db.DB, accessTokenRepository *postgres.AccessToken, userID int32) *entity.AccessToken {
					token := &entity.AccessToken{
						UserID:    userID,
						TokenHash: "test_exist",
						ExpiredAt: now.Add(time.Minute * 2),
						CreatedAt: now,
						UpdatedAt: now.Add(time.Minute),
					}

					accessTokenRepository.SetTimeNow(func() time.Time {
						return now
					})
					_, err := accessTokenRepository.CreateOrUpdate(ctx, userID, "test", now)
					if err != nil {
						t.Fatalf("CreateOrUpdate() err = %v; want nil", err)
					}

					accessTokenRepository.SetTimeNow(func() time.Time {
						return now.Add(time.Minute)
					})
					return token
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				accessTokenRepository := postgres.NewAccessToken(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				want := tc.setup(dbConn, accessTokenRepository, user.ID)

				got, err := accessTokenRepository.CreateOrUpdate(ctx, user.ID, want.TokenHash, want.ExpiredAt)
				if err != nil {
					t.Fatalf("accessToken.CreateOrUpdate() err = %v; want nil", err)
				}
				want.ID = got.ID
				compareAccessToken(t, got, want)

				got, err = accessTokenRepository.GetByTokenHash(ctx, want.TokenHash)
				if err != nil {
					t.Fatalf("GetByTokenHash() err = %v; want nil", err)
				}
				compareAccessToken(t, got, want)
			})
		}
	}
}
func testRepositoryAccessToken_DeleteByID(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(dbConn *db.DB, accessTokenRepository *postgres.AccessToken, userID int32) (int32, string, *entity.AccessToken)
		}{
			"delete exist record": {
				setup: func(dbConn *db.DB, accessTokenRepository *postgres.AccessToken, userID int32) (int32, string, *entity.AccessToken) {
					accessTokenRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := accessTokenRepository.CreateOrUpdate(ctx, userID, "test", now)
					if err != nil {
						t.Fatalf("CreateOrUpdate() err = %v; want nil", err)
					}

					return resp.ID, "test", nil
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				accessTokenRepository := postgres.NewAccessToken(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				tokenID, tokenHash, want := tc.setup(dbConn, accessTokenRepository, user.ID)

				err = accessTokenRepository.DeleteByID(ctx, tokenID)
				if err != nil {
					t.Fatalf("DeleteByID() err = %v; want nil", err)
				}

				got, err := accessTokenRepository.GetByTokenHash(ctx, tokenHash)
				if err != nil {
					t.Fatalf("GetByTokenHash() err = %v; want nil", err)
				}
				compareAccessToken(t, got, want)
			})
		}
	}
}
