package postgres

import (
	"context"
	"time"

	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
)

type IAccessToken interface {
	CreateOrUpdate(ctx context.Context, userID int32, tokenHash string, expiredAt time.Time) (*entity.AccessToken, error)
	GetByUserID(ctx context.Context, tokenHash string) (*entity.AccessToken, error)
	DeleteByID(ctx context.Context, id int32) error
}

type IApplication interface {
	Create(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error)
	List(ctx context.Context, req *dto.ApplicationListRequest) ([]*entity.ApplicationPublic, error)
	Item(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPublic, error)
	PrivateItem(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPrivate, error)
}

type IEmailValidation interface {
	CreateOrUpdate(ctx context.Context, userID int32, code string, expiredAt time.Time) (*entity.EmailValidation, error)
	GetByCode(ctx context.Context, code string) (*entity.EmailValidation, error)
	DeleteByID(ctx context.Context, id int32) error
}

type IHistory interface {
	NewEvent(ctx context.Context, userID int32, message string) error
	NewEventWithAffected(ctx context.Context, userID, affectedUserID int32, message string) error
}

type IPassword interface {
	Create(ctx context.Context, userID int32, passwordHash string) (*entity.Password, error)
	GetByUserID(ctx context.Context, userID int32) (*entity.Password, error)
	Update(ctx context.Context, id int32, passwordHash string) (*entity.Password, error)
	IncreaseFailedAttempts(ctx context.Context, id int32) (*entity.Password, error)
	ResetFailedAttempts(ctx context.Context, id int32) (*entity.Password, error)
}

type IUser interface {
	GetByID(ctx context.Context, id int32) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	Create(ctx context.Context, email, name string) (*entity.User, error)
	ActivateRecord(ctx context.Context, userID int32) (*entity.User, error)
	UpdateSteamID(ctx context.Context, userID int32, steamID string) (*entity.User, error)
}

type Postgres struct {
	txManager *txManager

	AccessToken     IAccessToken
	Application     IApplication
	EmailValidation IEmailValidation
	History         IHistory
	Password        IPassword
	User            IUser
}

func NewPostgres(db *db.DB,
	accessToken IAccessToken,
	application IApplication,
	emailValidation IEmailValidation,
	history IHistory,
	password IPassword,
	user IUser,
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

func (r *Postgres) TxManager() *txManager {
	return r.txManager
}
