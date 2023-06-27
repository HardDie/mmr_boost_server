package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/utils"
)

type User struct {
	repository *postgres.Postgres
}

func NewUser(repository *postgres.Postgres) *User {
	return &User{
		repository: repository,
	}
}

func (s *User) UpdatePassword(ctx context.Context, req *dto.UserUpdatePasswordRequest, userID int32) error {
	err := s.repository.TxManager().ReadWriteTx(ctx, func(ctx context.Context) error {
		// Get password from DB
		password, err := s.repository.Password.GetByUserID(ctx, userID)
		if err != nil {
			logger.Error.Printf("error read password from DB: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		if password == nil {
			logger.Error.Printf("password for user %d not found in DB", userID)
			return status.Error(codes.Internal, "internal")
		}

		// Check if password is correct
		if !utils.HashBcryptCompare(req.OldPassword, password.PasswordHash) {
			return status.Error(codes.InvalidArgument, "invalid old password")
		}

		// Hashing password
		hashPassword, err := utils.HashBcrypt(req.NewPassword)
		if err != nil {
			logger.Error.Printf("error hashing password: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}

		// Update password
		_, err = s.repository.Password.Update(ctx, userID, hashPassword)
		if err != nil {
			logger.Error.Printf("error updating password in DB: %v", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *User) UpdateSteamID(ctx context.Context, req *dto.UserUpdateSteamIDRequest, userID int32) (*entity.User, error) {
	u, err := s.repository.User.UpdateSteamID(ctx, userID, req.SteamID)
	if err != nil {
		return nil, err
	}
	return u, nil
}
