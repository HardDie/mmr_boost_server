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
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

func compareUser(t *testing.T, a, b *entity.User) {
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
	if a.Email != b.Email {
		t.Errorf("Email = %q; want %q", a.Email, b.Email)
	}
	if a.Username != b.Username {
		t.Errorf("Username = %q; want %q", a.Username, b.Username)
	}
	if a.RoleID != b.RoleID {
		t.Errorf("RoleID = %d; want %d", a.RoleID, b.RoleID)
	}
	if !utils.Compare(a.SteamID, b.SteamID) {
		t.Errorf("SteamID = %v; want %v", a.SteamID, b.SteamID)
	}
	if a.IsActivated != b.IsActivated {
		t.Errorf("IsActivated = %v; want %v", a.IsActivated, b.IsActivated)
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

func TestRepositoryUser(t *testing.T) {
	dbConn := setup(t)

	t.Run("GetByID", testRepositoryUser_GetByID(dbConn))
	t.Run("GetByName", testRepositoryUser_GetByName(dbConn))
	t.Run("GetByNameOrEmail", testRepositoryUser_GetByNameOrEmail(dbConn))
	t.Run("Create", testRepositoryUser_Create(dbConn))
	t.Run("ActivateRecord", testRepositoryUser_ActivateRecord(dbConn))
	t.Run("UpdateSteamID", testRepositoryUser_UpdateSteamID(dbConn))
}

func testRepositoryUser_GetByID(dbConn *db.DB) func(*testing.T) {
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(t *testing.T, db *db.DB) (int32, *entity.User)
		}{
			"manually created user": {
				setup: func(t *testing.T, db *db.DB) (int32, *entity.User) {
					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now,
					}

					q := gosql.NewInsert().Into("users")
					q.Columns().Add("email", "username", "role_id", "is_activated", "created_at", "updated_at")
					q.Columns().Arg(user.Email, user.Username, user.RoleID, user.IsActivated, user.CreatedAt, user.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&user.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}
					return user.ID, user
				},
			},
			"not exist user": {
				setup: func(t *testing.T, db *db.DB) (int32, *entity.User) {
					return 10, nil
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userID, u := tc.setup(t, dbConn)

				userRepository := postgres.NewUser(dbConn)
				resp, err := userRepository.GetByID(ctx, userID)
				if err != nil {
					t.Fatalf("GetByID() err = %v; want nil", err)
				}
				if resp == nil && resp != u {
					t.Fatalf("GetByID() = nil; want %v", u)
				}
				compareUser(t, resp, u)
			})
		}
	}
}
func testRepositoryUser_GetByName(dbConn *db.DB) func(*testing.T) {
	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(t *testing.T, db *db.DB) (string, *entity.User)
		}{
			"manually created user": {
				setup: func(t *testing.T, db *db.DB) (string, *entity.User) {
					now := utils.TimeTrim(time.Now())

					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now,
					}

					q := gosql.NewInsert().Into("users")
					q.Columns().Add("email", "username", "role_id", "is_activated", "created_at", "updated_at")
					q.Columns().Arg(user.Email, user.Username, user.RoleID, user.IsActivated, user.CreatedAt, user.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&user.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}
					return user.Username, user
				},
			},
			"not exist user": {
				setup: func(t *testing.T, db *db.DB) (string, *entity.User) {
					return "unknown", nil
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userName, u := tc.setup(t, dbConn)

				userRepository := postgres.NewUser(dbConn)
				resp, err := userRepository.GetByName(ctx, userName)
				if err != nil {
					t.Fatalf("GetByName() err = %v; want nil", err)
				}
				if resp == nil && resp != u {
					t.Fatalf("GetByName() = nil; want %v", u)
				}
				compareUser(t, resp, u)
			})
		}
	}
}
func testRepositoryUser_GetByNameOrEmail(dbConn *db.DB) func(*testing.T) {
	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(t *testing.T, db *db.DB) (string, string, *entity.User)
		}{
			"correct name and email": {
				setup: func(t *testing.T, db *db.DB) (string, string, *entity.User) {
					now := time.Now()
					now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now,
					}

					q := gosql.NewInsert().Into("users")
					q.Columns().Add("email", "username", "role_id", "is_activated", "created_at", "updated_at")
					q.Columns().Arg(user.Email, user.Username, user.RoleID, user.IsActivated, user.CreatedAt, user.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&user.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}
					return user.Username, user.Email, user
				},
			},
			"correct name": {
				setup: func(t *testing.T, db *db.DB) (string, string, *entity.User) {
					now := time.Now()
					now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now,
					}

					q := gosql.NewInsert().Into("users")
					q.Columns().Add("email", "username", "role_id", "is_activated", "created_at", "updated_at")
					q.Columns().Arg(user.Email, user.Username, user.RoleID, user.IsActivated, user.CreatedAt, user.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&user.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}
					return user.Username, "invalid@mail.com", user
				},
			},
			"correct email": {
				setup: func(t *testing.T, db *db.DB) (string, string, *entity.User) {
					now := time.Now()
					now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now,
					}

					q := gosql.NewInsert().Into("users")
					q.Columns().Add("email", "username", "role_id", "is_activated", "created_at", "updated_at")
					q.Columns().Arg(user.Email, user.Username, user.RoleID, user.IsActivated, user.CreatedAt, user.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&user.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}
					return "unknown", user.Email, user
				},
			},
			"not exist user": {
				setup: func(t *testing.T, db *db.DB) (string, string, *entity.User) {
					return "unknown", "invalid@mail.com", nil
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userName, userEmail, u := tc.setup(t, dbConn)

				userRepository := postgres.NewUser(dbConn)
				resp, err := userRepository.GetByNameOrEmail(ctx, userName, userEmail)
				if err != nil {
					t.Fatalf("GetByNameOrEmail() err = %v; want nil", err)
				}
				if resp == nil && resp != u {
					t.Fatalf("GetByNameOrEmail() = nil; want %v", u)
				}
				compareUser(t, resp, u)
			})
		}
	}
}
func testRepositoryUser_Create(dbConn *db.DB) func(*testing.T) {
	return func(t *testing.T) {
		now := utils.TimeTrim(time.Now())

		tests := map[string]struct {
			req *entity.User
		}{
			"valid": {
				req: &entity.User{
					Email:       "test@mail.com",
					Username:    "test",
					RoleID:      int32(pb.UserRoleID_user),
					IsActivated: false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				userRepository.SetTimeNow(func() time.Time {
					return now
				})

				resp, err := userRepository.Create(ctx, tc.req.Email, tc.req.Username)
				if err != nil {
					t.Fatalf("Create() err = %v; want nil", err)
				}
				// Update ID
				tc.req.ID = resp.ID
				compareUser(t, resp, tc.req)

				resp, err = userRepository.GetByID(ctx, resp.ID)
				if err != nil {
					t.Fatalf("GetByID() err = %v; want nil", err)
				}
				compareUser(t, resp, tc.req)
			})
		}
	}
}
func testRepositoryUser_ActivateRecord(dbConn *db.DB) func(*testing.T) {
	return func(t *testing.T) {
		now := utils.TimeTrim(time.Now())
		ctx := context.Background()

		tests := map[string]struct {
			setup func(t *testing.T, userRepository *postgres.User) *entity.User
		}{
			"activate new created record": {
				setup: func(t *testing.T, userRepository *postgres.User) *entity.User {
					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now.Add(time.Minute),
					}

					userRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := userRepository.Create(ctx, user.Email, user.Username)
					if err != nil {
						t.Fatalf("Create() err = %v; want nil", err)
					}
					user.ID = resp.ID

					return user
				},
			},
			"error user already activated": {
				setup: func(t *testing.T, userRepository *postgres.User) *entity.User {
					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now.Add(time.Minute),
					}

					userRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := userRepository.Create(ctx, user.Email, user.Username)
					if err != nil {
						t.Fatalf("Create() err = %v; want nil", err)
					}
					user.ID = resp.ID

					userRepository.SetTimeNow(func() time.Time {
						return now.Add(time.Minute)
					})
					_, err = userRepository.ActivateRecord(ctx, user.ID)
					if err != nil {
						t.Fatalf("ActivateRecord() err = %v; want nil", err)
					}

					return user
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				want := tc.setup(t, userRepository)

				userRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})
				got, err := userRepository.ActivateRecord(ctx, want.ID)
				if err != nil {
					t.Fatalf("ActivateRecord() err = %v; want nil", err)
				}
				compareUser(t, got, want)
			})
		}
	}
}
func testRepositoryUser_UpdateSteamID(dbConn *db.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		now := utils.TimeTrim(time.Now())

		tests := map[string]struct {
			setup func(t *testing.T, userRepository *postgres.User) *entity.User
		}{
			"set steam id for activated record": {
				setup: func(t *testing.T, userRepository *postgres.User) *entity.User {
					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						SteamID:     utils.Allocate("someTestSteamID"),
						IsActivated: true,
						CreatedAt:   now,
						UpdatedAt:   now.Add(time.Minute),
					}

					userRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := userRepository.Create(ctx, user.Email, user.Username)
					if err != nil {
						t.Fatalf("Create() err = %v; want nil", err)
					}
					user.ID = resp.ID

					userRepository.SetTimeNow(func() time.Time {
						return now.Add(time.Minute)
					})
					_, err = userRepository.ActivateRecord(ctx, user.ID)
					if err != nil {
						t.Fatalf("ActivateRecord() err = %v; want nil", err)
					}

					return user
				},
			},
			"set steam id for new record": {
				setup: func(t *testing.T, userRepository *postgres.User) *entity.User {
					user := &entity.User{
						Email:       "test@mail.com",
						Username:    "test",
						RoleID:      int32(pb.UserRoleID_user),
						SteamID:     utils.Allocate("someTestSteamID"),
						IsActivated: false,
						CreatedAt:   now,
						UpdatedAt:   now.Add(time.Minute),
					}

					userRepository.SetTimeNow(func() time.Time {
						return now
					})
					resp, err := userRepository.Create(ctx, user.Email, user.Username)
					if err != nil {
						t.Fatalf("Create() err = %v; want nil", err)
					}
					user.ID = resp.ID

					return user
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)

				want := tc.setup(t, userRepository)

				userRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})
				got, err := userRepository.UpdateSteamID(ctx, want.ID, *want.SteamID)
				if err != nil {
					t.Fatalf("UpdateSteamID() err = %v; want nil", err)
				}
				compareUser(t, got, want)

				got, err = userRepository.GetByID(ctx, want.ID)
				if err != nil {
					t.Fatalf("GetByID() err = %v; want nil", err)
				}
				compareUser(t, got, want)
			})
		}
	}
}
