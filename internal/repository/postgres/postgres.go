package postgres

import (
	"context"
	"time"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
)

type IPostgresAccessToken interface {
	CreateOrUpdate(ctx context.Context, userID int32, tokenHash string, expiredAt time.Time) (*entity.AccessToken, error)
	GetByUserID(ctx context.Context, tokenHash string) (*entity.AccessToken, error)
	DeleteByID(ctx context.Context, id int32) error
}

type IPostgresApplication interface {
	Create(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error)
	List(ctx context.Context, req *dto.ApplicationListRequest) ([]*entity.ApplicationPublic, error)
	Item(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPublic, error)
	PrivateItem(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPrivate, error)
	UpdateStatus(ctx context.Context, req *dto.ApplicationUpdateStatusRequest) (*entity.ApplicationPublic, error)
	UpdateItem(ctx context.Context, req *dto.ApplicationUpdateRequest) (*entity.ApplicationPublic, error)
	UpdatePrivate(ctx context.Context, req *dto.ApplicationUpdatePrivateRequest) (*entity.ApplicationPrivate, error)
}

type IPostgresEmailValidation interface {
	CreateOrUpdate(ctx context.Context, userID int32, code string, expiredAt time.Time) (*entity.EmailValidation, error)
	GetByCode(ctx context.Context, code string) (*entity.EmailValidation, error)
	DeleteByID(ctx context.Context, id int32) error
}

type IPostgresHistory interface {
	NewEvent(ctx context.Context, userID int32, message string) error
	NewEventWithAffected(ctx context.Context, userID, affectedUserID int32, message string) error
}

type IPostgresPassword interface {
	Create(ctx context.Context, userID int32, passwordHash string) (*entity.Password, error)
	GetByUserID(ctx context.Context, userID int32) (*entity.Password, error)
	Update(ctx context.Context, id int32, passwordHash string) (*entity.Password, error)
	IncreaseFailedAttempts(ctx context.Context, id int32) (*entity.Password, error)
	ResetFailedAttempts(ctx context.Context, id int32) (*entity.Password, error)
}

type IPostgresUser interface {
	GetByID(ctx context.Context, id int32) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	GetByNameOrEmail(ctx context.Context, name string, email string) (*entity.User, error)
	Create(ctx context.Context, email, name string) (*entity.User, error)
	ActivateRecord(ctx context.Context, userID int32) (*entity.User, error)
	UpdateSteamID(ctx context.Context, userID int32, steamID string) (*entity.User, error)
}

type Postgres struct {
	txManager *txManager

	AccessToken     IPostgresAccessToken
	Application     IPostgresApplication
	EmailValidation IPostgresEmailValidation
	History         IPostgresHistory
	Password        IPostgresPassword
	User            IPostgresUser
}

func NewPostgres(db *db.DB,
	accessToken IPostgresAccessToken,
	application IPostgresApplication,
	emailValidation IPostgresEmailValidation,
	history IPostgresHistory,
	password IPostgresPassword,
	user IPostgresUser,
) *Postgres {
	return &Postgres{
		txManager: newTxManager(db),

		AccessToken:     accessToken,
		Application:     application,
		EmailValidation: emailValidation,
		History:         history,
		Password:        password,
		User:            user,
	}
}

func (r *Postgres) TxManager() *txManager { //nolint:revive
	return r.txManager
}
