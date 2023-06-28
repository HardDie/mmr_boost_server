package service

import (
	"context"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
)

type IServiceApplication interface {
	Create(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error)
	UserList(ctx context.Context, req *dto.ApplicationUserListRequest) ([]*entity.ApplicationPublic, error)
	ManagementList(ctx context.Context, req *dto.ApplicationManagementListRequest) ([]*entity.ApplicationPublic, error)
	UserItem(ctx context.Context, req *dto.ApplicationUserItemRequest) (*entity.ApplicationPublic, error)
	DeleteItem(ctx context.Context, req *dto.ApplicationItemDeleteRequest) error

	ManagementItem(ctx context.Context, req *dto.ApplicationManagementItemRequest) (*entity.ApplicationPublic, error)
	ManagementPrivateItem(ctx context.Context, req *dto.ApplicationManagementItemRequest) (*entity.ApplicationPrivate, error)
}

type IServiceAuth interface {
	Register(ctx context.Context, req *dto.AuthRegisterRequest) error
	Login(ctx context.Context, req *dto.AuthLoginRequest) (*entity.User, error)
	Logout(ctx context.Context, sessionID int32) error
	GenerateCookie(ctx context.Context, userID int32) (*entity.AccessToken, error)
	ValidateCookie(ctx context.Context, sessionKey string) (*entity.User, *entity.AccessToken, error)
	GetUserInfo(ctx context.Context, userID int32) (*entity.User, error)
	ValidateEmail(ctx context.Context, code string) error
	SendValidationEmail(ctx context.Context, name string) error
}

type IServiceSystem interface {
	GetSwagger() ([]byte, error)
}

type IServiceUser interface {
	UpdatePassword(ctx context.Context, req *dto.UserUpdatePasswordRequest, userID int32) error
	UpdateSteamID(ctx context.Context, req *dto.UserUpdateSteamIDRequest, userID int32) (*entity.User, error)
}

type IServicePrice interface {
	Price(ctx context.Context, req *dto.PriceRequest) (float64, error)
}

type Service struct {
	Application IServiceApplication
	Auth        IServiceAuth
	System      IServiceSystem
	User        IServiceUser
	Price       IServicePrice
}

func NewService(
	application IServiceApplication,
	auth IServiceAuth,
	system IServiceSystem,
	user IServiceUser,
	price IServicePrice,
) *Service {
	return &Service{
		Application: application,
		Auth:        auth,
		System:      system,
		User:        user,
		Price:       price,
	}
}
