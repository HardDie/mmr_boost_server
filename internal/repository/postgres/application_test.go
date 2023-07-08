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
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

func compareApplicationPublic(t *testing.T, a, b *entity.ApplicationPublic) {
	t.Helper()

	if a == nil && b == nil {
		return
	}
	if a == nil && b != nil {
		t.Fatalf("ApplicationPublic = nil; want not nil")
	}
	if a != nil && b == nil {
		t.Fatalf("ApplicationPublic = not nil; want nil")
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
	if a.StatusID != b.StatusID {
		t.Errorf("StatusID = %d; want %d", a.StatusID, b.StatusID)
	}
	if a.TypeID != b.TypeID {
		t.Errorf("TypeID = %d; want %d", a.TypeID, b.TypeID)
	}
	if a.CurrentMMR != b.CurrentMMR {
		t.Errorf("CurrentMMR = %d; want %d", a.CurrentMMR, b.CurrentMMR)
	}
	if a.TargetMMR != b.TargetMMR {
		t.Errorf("TargetMMR = %d; want %d", a.TargetMMR, b.TargetMMR)
	}
	if a.TgContact != b.TgContact {
		t.Errorf("TgContact = %q; want %q", a.TgContact, b.TgContact)
	}
	if a.Price != b.Price {
		t.Errorf("Price = %d; want %d", a.Price, b.Price)
	}
	if a.Comment != b.Comment {
		t.Errorf("Comment = %q; want %q", a.Comment, b.Comment)
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
func compareApplicationPrivate(t *testing.T, a, b *entity.ApplicationPrivate) {
	t.Helper()

	if a == nil && b == nil {
		return
	}
	if a == nil && b != nil {
		t.Fatalf("ApplicationPublic = nil; want not nil")
	}
	if a != nil && b == nil {
		t.Fatalf("ApplicationPublic = not nil; want nil")
	}
	if reflect.DeepEqual(a, b) {
		return
	}
	if a.ID != b.ID {
		t.Errorf("ID = %d; want %d", a.ID, b.ID)
	}
	if !utils.Compare(a.SteamLogin, b.SteamLogin) {
		t.Errorf("SteamLogin = %v; want %v", a.SteamLogin, b.SteamLogin)
	}
	if !utils.Compare(a.SteamPassword, b.SteamPassword) {
		t.Errorf("SteamPassword = %v; want %v", a.SteamPassword, b.SteamPassword)
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

func TestRepositoryApplication(t *testing.T) {
	dbConn := setup(t)

	t.Run("Item", testRepositoryApplication_Item(dbConn))
	t.Run("Create", testRepositoryApplication_Create(dbConn))
	t.Run("List", testRepositoryApplication_List(dbConn))
	t.Run("PrivateItem", testRepositoryApplication_PrivateItem(dbConn))
	t.Run("UpdateStatus", testRepositoryApplication_UpdateStatus(dbConn))
	t.Run("UpdatePrivate", testRepositoryApplication_UpdatePrivate(dbConn))
	t.Run("UpdateItem", testRepositoryApplication_UpdateItem(dbConn))
}

func testRepositoryApplication_Item(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(userID int32, db *db.DB) (*dto.ApplicationItemRequest, *entity.ApplicationPublic)
		}{
			"manually created application get by id": {
				setup: func(userID int32, db *db.DB) (*dto.ApplicationItemRequest, *entity.ApplicationPublic) {
					application := &entity.ApplicationPublic{
						UserID:     userID,
						StatusID:   int32(pb.ApplicationStatusID_created),
						TypeID:     int32(pb.ApplicationTypeID_boost_mmr),
						CurrentMMR: 1000,
						TargetMMR:  2000,
						TgContact:  "vas100",
						Price:      800,
						CreatedAt:  now,
						UpdatedAt:  now,
					}

					q := gosql.NewInsert().Into("applications")
					q.Columns().Add("user_id", "status_id", "type_id", "current_mmr", "target_mmr", "tg_contact", "price", "created_at", "updated_at")
					q.Columns().Arg(application.UserID, application.StatusID, application.TypeID,
						application.CurrentMMR, application.TargetMMR, application.TgContact, application.Price,
						application.CreatedAt, application.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&application.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}

					return &dto.ApplicationItemRequest{ApplicationID: application.ID}, application
				},
			},
			"manually created application get by user_id": {
				setup: func(userID int32, db *db.DB) (*dto.ApplicationItemRequest, *entity.ApplicationPublic) {
					application := &entity.ApplicationPublic{
						UserID:     userID,
						StatusID:   int32(pb.ApplicationStatusID_created),
						TypeID:     int32(pb.ApplicationTypeID_boost_mmr),
						CurrentMMR: 1000,
						TargetMMR:  2000,
						TgContact:  "vas100",
						Price:      800,
						CreatedAt:  now,
						UpdatedAt:  now,
					}

					q := gosql.NewInsert().Into("applications")
					q.Columns().Add("user_id", "status_id", "type_id", "current_mmr", "target_mmr", "tg_contact", "price", "created_at", "updated_at")
					q.Columns().Arg(application.UserID, application.StatusID, application.TypeID,
						application.CurrentMMR, application.TargetMMR, application.TgContact, application.Price,
						application.CreatedAt, application.UpdatedAt)
					q.Returning().Add("id")

					err := db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&application.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}

					return &dto.ApplicationItemRequest{ApplicationID: application.ID, UserID: &application.UserID}, application
				},
			},
			"not exist application": {
				setup: func(userID int32, db *db.DB) (*dto.ApplicationItemRequest, *entity.ApplicationPublic) {
					return &dto.ApplicationItemRequest{ApplicationID: 10}, nil
				},
			},
		}

		ctx := context.Background()
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				applicationRepository := postgres.NewApplication(dbConn)

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				arg, want := tc.setup(user.ID, dbConn)

				got, err := applicationRepository.Item(ctx, arg)
				if err != nil {
					t.Fatalf("GetByTokenHash() err = %v; want nil", err)
				}
				compareApplicationPublic(t, got, want)
			})
		}
	}
}
func testRepositoryApplication_Create(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(userID int32) (*dto.ApplicationCreateRequest, *entity.ApplicationPublic)
		}{
			"valid boost mmr": {
				setup: func(userID int32) (*dto.ApplicationCreateRequest, *entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:     userID,
						TypeID:     int32(pb.ApplicationTypeID_boost_mmr),
						CurrentMMR: 1000,
						TargetMMR:  2000,
						TgContact:  "vasiliy",
						Price:      200,
					}
					want := &entity.ApplicationPublic{
						UserID:     arg.UserID,
						StatusID:   int32(pb.ApplicationStatusID_created),
						TypeID:     arg.TypeID,
						CurrentMMR: arg.CurrentMMR,
						TargetMMR:  arg.TargetMMR,
						TgContact:  arg.TgContact,
						Price:      arg.Price,
						CreatedAt:  now,
						UpdatedAt:  now,
					}
					return arg, want
				},
			},
			"valid calibration": {
				setup: func(userID int32) (*dto.ApplicationCreateRequest, *entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					want := &entity.ApplicationPublic{
						UserID:    arg.UserID,
						StatusID:  int32(pb.ApplicationStatusID_created),
						TypeID:    arg.TypeID,
						TgContact: arg.TgContact,
						Price:     arg.Price,
						CreatedAt: now,
						UpdatedAt: now,
					}
					return arg, want
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				applicationRepository := postgres.NewApplication(dbConn)
				applicationRepository.SetTimeNow(func() time.Time {
					return now
				})

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				arg, want := tc.setup(user.ID)

				got, err := applicationRepository.Create(ctx, arg)
				if err != nil {
					t.Fatalf("Create() err = %v; want nil", err)
				}
				want.ID = got.ID
				compareApplicationPublic(t, got, want)

				got, err = applicationRepository.Item(ctx, &dto.ApplicationItemRequest{
					ApplicationID: want.ID,
				})
				if err != nil {
					t.Fatalf("Item() err = %v; want nil", err)
				}
				compareApplicationPublic(t, got, want)
			})
		}
	}
}
func testRepositoryApplication_List(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic)
		}{
			"empty list": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					return &dto.ApplicationListRequest{}, []*entity.ApplicationPublic{}
				},
			},
			"single record": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{}, []*entity.ApplicationPublic{resp}
				},
			},
			"two record": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					resp2, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{}, []*entity.ApplicationPublic{resp2, resp}
				},
			},
			"filter by status created": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{StatusID: utils.Allocate(int32(pb.ApplicationStatusID_created))},
						[]*entity.ApplicationPublic{resp}
				},
			},
			"filter by status paid": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					_, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{StatusID: utils.Allocate(int32(pb.ApplicationStatusID_paid))},
						[]*entity.ApplicationPublic{}
				},
			},
			"filter by user valid": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{UserID: utils.Allocate(userID)},
						[]*entity.ApplicationPublic{resp}
				},
			},
			"filter by user invalid": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					_, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{UserID: utils.Allocate(userID + 1)},
						[]*entity.ApplicationPublic{}
				},
			},
			"filter by status and user valid": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{UserID: utils.Allocate(userID), StatusID: utils.Allocate(int32(pb.ApplicationStatusID_created))},
						[]*entity.ApplicationPublic{resp}
				},
			},
			"filter by status and user invalid": {
				setup: func(applicationRepository *postgres.Application, userID int32) (*dto.ApplicationListRequest, []*entity.ApplicationPublic) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					_, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationListRequest{UserID: utils.Allocate(userID + 1), StatusID: utils.Allocate(int32(pb.ApplicationStatusID_paid))},
						[]*entity.ApplicationPublic{}
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				applicationRepository := postgres.NewApplication(dbConn)
				applicationRepository.SetTimeNow(func() time.Time {
					return now
				})

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				arg, want := tc.setup(applicationRepository, user.ID)

				got, err := applicationRepository.List(ctx, arg)
				if err != nil {
					t.Fatalf("Create() err = %v; want nil", err)
				}
				if len(got) != len(want) {
					t.Fatalf("len(got) = %d; want %d", len(got), len(want))
				}
				for i := range got {
					compareApplicationPublic(t, got[i], want[i])
				}
			})
		}
	}
}
func testRepositoryApplication_PrivateItem(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(db *db.DB, applicationRepository *postgres.Application, userID int32) (*dto.ApplicationItemRequest, *entity.ApplicationPrivate)
		}{
			"empty record": {
				setup: func(db *db.DB, applicationRepository *postgres.Application, userID int32) (*dto.ApplicationItemRequest, *entity.ApplicationPrivate) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					want := &entity.ApplicationPrivate{
						ID:        resp.ID,
						CreatedAt: now,
						UpdatedAt: now,
					}
					return &dto.ApplicationItemRequest{ApplicationID: resp.ID}, want
				},
			},
			"manually filled record": {
				setup: func(db *db.DB, applicationRepository *postgres.Application, userID int32) (*dto.ApplicationItemRequest, *entity.ApplicationPrivate) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}

					q := gosql.NewUpdate().Table("applications")
					q.Set().Append("steam_login = ?", "test1")
					q.Set().Append("steam_password = ?", "test2")
					q.Set().Append("updated_at = ?", now.Add(time.Minute))
					q.Where().AddExpression("id = ?", resp.ID)
					q.Returning().Add("id")

					err = db.DB.QueryRow(q.String(), q.GetArguments()...).Scan(&resp.ID)
					if err != nil {
						t.Fatalf("QueryRow() err = %v; want nil", err)
					}

					want := &entity.ApplicationPrivate{
						ID:            resp.ID,
						SteamLogin:    utils.Allocate("test1"),
						SteamPassword: utils.Allocate("test2"),
						CreatedAt:     now,
						UpdatedAt:     now.Add(time.Minute),
					}
					return &dto.ApplicationItemRequest{ApplicationID: resp.ID}, want
				},
			},
			"invalid record": {
				setup: func(db *db.DB, applicationRepository *postgres.Application, userID int32) (*dto.ApplicationItemRequest, *entity.ApplicationPrivate) {
					arg := &dto.ApplicationCreateRequest{
						UserID:    userID,
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						TgContact: "vasiliy",
						Price:     2000,
					}
					resp, err := applicationRepository.Create(ctx, arg)
					if err != nil {
						t.Fatalf("application.Create() err = %v; want nil", err)
					}
					return &dto.ApplicationItemRequest{ApplicationID: resp.ID + 1}, nil
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				applicationRepository := postgres.NewApplication(dbConn)
				applicationRepository.SetTimeNow(func() time.Time {
					return now
				})

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}

				arg, want := tc.setup(dbConn, applicationRepository, user.ID)

				got, err := applicationRepository.PrivateItem(ctx, arg)
				if err != nil {
					t.Fatalf("PrivateItem() err = %v; want nil", err)
				}
				compareApplicationPrivate(t, got, want)
			})
		}
	}
}
func testRepositoryApplication_UpdateStatus(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(userID, applicationID int32) (int32, *entity.ApplicationPublic)
		}{
			"update on new status": {
				setup: func(userID, applicationID int32) (int32, *entity.ApplicationPublic) {
					return int32(pb.ApplicationStatusID_awaits_payment), &entity.ApplicationPublic{
						ID:        applicationID,
						UserID:    userID,
						StatusID:  int32(pb.ApplicationStatusID_awaits_payment),
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						CreatedAt: now,
						UpdatedAt: now.Add(time.Minute),
					}
				},
			},
			"update on same status": {
				setup: func(userID, applicationID int32) (int32, *entity.ApplicationPublic) {
					return int32(pb.ApplicationStatusID_created), &entity.ApplicationPublic{
						ID:        applicationID,
						UserID:    userID,
						StatusID:  int32(pb.ApplicationStatusID_created),
						TypeID:    int32(pb.ApplicationTypeID_calibration),
						CreatedAt: now,
						UpdatedAt: now.Add(time.Minute),
					}
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				applicationRepository := postgres.NewApplication(dbConn)
				applicationRepository.SetTimeNow(func() time.Time {
					return now
				})

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}
				application, err := applicationRepository.Create(ctx, &dto.ApplicationCreateRequest{
					UserID: user.ID,
					TypeID: int32(pb.ApplicationTypeID_calibration),
				})
				if err != nil {
					t.Fatalf("application.Create() err = %v; want nil", err)
				}

				arg, want := tc.setup(user.ID, application.ID)

				applicationRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})

				got, err := applicationRepository.UpdateStatus(ctx, &dto.ApplicationUpdateStatusRequest{ApplicationID: application.ID, StatusID: arg})
				if err != nil {
					t.Fatalf("UpdateStatus() err = %v; want nil", err)
				}
				compareApplicationPublic(t, got, want)

				got, err = applicationRepository.Item(ctx, &dto.ApplicationItemRequest{ApplicationID: application.ID})
				if err != nil {
					t.Fatalf("Item() err = %v; want nil", err)
				}
				compareApplicationPublic(t, got, want)
			})
		}
	}
}
func testRepositoryApplication_UpdatePrivate(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(userID, applicationID int32) (*dto.ApplicationUpdatePrivateRequest, *entity.ApplicationPrivate)
		}{
			"update on private info": {
				setup: func(userID, applicationID int32) (*dto.ApplicationUpdatePrivateRequest, *entity.ApplicationPrivate) {
					arg := &dto.ApplicationUpdatePrivateRequest{
						ApplicationID: applicationID,
						SteamLogin:    "test1",
						SteamPassword: "test2",
					}
					want := &entity.ApplicationPrivate{
						ID:            applicationID,
						SteamLogin:    utils.Allocate("test1"),
						SteamPassword: utils.Allocate("test2"),
						CreatedAt:     now,
						UpdatedAt:     now.Add(time.Minute),
					}
					return arg, want
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				applicationRepository := postgres.NewApplication(dbConn)
				applicationRepository.SetTimeNow(func() time.Time {
					return now
				})

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}
				application, err := applicationRepository.Create(ctx, &dto.ApplicationCreateRequest{
					UserID: user.ID,
					TypeID: int32(pb.ApplicationTypeID_calibration),
				})
				if err != nil {
					t.Fatalf("application.Create() err = %v; want nil", err)
				}

				arg, want := tc.setup(user.ID, application.ID)

				applicationRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})

				got, err := applicationRepository.UpdatePrivate(ctx, arg)
				if err != nil {
					t.Fatalf("UpdatePrivate() err = %v; want nil", err)
				}
				compareApplicationPrivate(t, got, want)

				got, err = applicationRepository.PrivateItem(ctx, &dto.ApplicationItemRequest{ApplicationID: application.ID})
				if err != nil {
					t.Fatalf("PrivateItem() err = %v; want nil", err)
				}
				compareApplicationPrivate(t, got, want)
			})
		}
	}
}
func testRepositoryApplication_UpdateItem(dbConn *db.DB) func(t *testing.T) {
	now := utils.TimeTrim(time.Now())
	ctx := context.Background()

	return func(t *testing.T) {
		tests := map[string]struct {
			setup func(application *entity.ApplicationPublic) (*dto.ApplicationUpdateRequest, *entity.ApplicationPublic)
		}{
			"valid": {
				setup: func(application *entity.ApplicationPublic) (*dto.ApplicationUpdateRequest, *entity.ApplicationPublic) {
					arg := &dto.ApplicationUpdateRequest{
						ApplicationID: application.ID,
						CurrentMMR:    2000,
						TargetMMR:     3000,
						TgContact:     utils.Allocate("newtg"),
						Price:         utils.Allocate[int32](7000),
					}
					want := &entity.ApplicationPublic{
						ID:         application.ID,
						UserID:     application.UserID,
						StatusID:   int32(pb.ApplicationStatusID_created),
						TypeID:     int32(pb.ApplicationTypeID_calibration),
						CurrentMMR: arg.CurrentMMR,
						TargetMMR:  arg.TargetMMR,
						TgContact:  *arg.TgContact,
						Price:      *arg.Price,
						CreatedAt:  now,
						UpdatedAt:  now.Add(time.Minute),
					}
					return arg, want
				},
			},
			"update to same values": {
				setup: func(application *entity.ApplicationPublic) (*dto.ApplicationUpdateRequest, *entity.ApplicationPublic) {
					arg := &dto.ApplicationUpdateRequest{
						ApplicationID: application.ID,
						CurrentMMR:    application.CurrentMMR,
						TargetMMR:     application.TargetMMR,
						TgContact:     utils.Allocate(application.TgContact),
						Price:         utils.Allocate[int32](application.Price),
					}
					want := &entity.ApplicationPublic{
						ID:         application.ID,
						UserID:     application.UserID,
						StatusID:   int32(pb.ApplicationStatusID_created),
						TypeID:     int32(pb.ApplicationTypeID_calibration),
						CurrentMMR: arg.CurrentMMR,
						TargetMMR:  arg.TargetMMR,
						TgContact:  *arg.TgContact,
						Price:      *arg.Price,
						CreatedAt:  now,
						UpdatedAt:  now.Add(time.Minute),
					}
					return arg, want
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				cleanDB(t, dbConn)
				userRepository := postgres.NewUser(dbConn)
				applicationRepository := postgres.NewApplication(dbConn)
				applicationRepository.SetTimeNow(func() time.Time {
					return now
				})

				user, err := userRepository.Create(ctx, "test@mail.com", "test")
				if err != nil {
					t.Fatalf("user.Create() err = %v; want nil", err)
				}
				application, err := applicationRepository.Create(ctx, &dto.ApplicationCreateRequest{
					UserID: user.ID,
					TypeID: int32(pb.ApplicationTypeID_calibration),
				})
				if err != nil {
					t.Fatalf("application.Create() err = %v; want nil", err)
				}

				arg, want := tc.setup(application)
				applicationRepository.SetTimeNow(func() time.Time {
					return now.Add(time.Minute)
				})

				got, err := applicationRepository.UpdateItem(ctx, arg)
				if err != nil {
					t.Fatalf("UpdateStatus() err = %v; want nil", err)
				}
				compareApplicationPublic(t, got, want)

				got, err = applicationRepository.Item(ctx, &dto.ApplicationItemRequest{ApplicationID: application.ID})
				if err != nil {
					t.Fatalf("Item() err = %v; want nil", err)
				}
				compareApplicationPublic(t, got, want)
			})
		}
	}
}
