package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type Application struct {
	repository *postgres.Postgres
}

func NewApplication(repository *postgres.Postgres) *Application {
	return &Application{
		repository: repository,
	}
}

func (s *Application) Create(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error) {
	var res *entity.ApplicationPublic

	err := s.repository.TxManager().ReadWriteTx(ctx, func(ctx context.Context) error {
		items, err := s.repository.Application.List(ctx, &dto.ApplicationListRequest{
			UserID:   &req.UserID,
			StatusID: utils.Allocate(int32(pb.ApplicationStatusID_created)),
		})
		if err != nil {
			logger.Error.Println("error get list of applications:", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		if len(items) != 0 {
			return status.Error(codes.InvalidArgument, "user already have active application")
		}

		resp, err := s.repository.Application.Create(ctx, req)
		if err != nil {
			logger.Error.Println("error creating new application:", err.Error())
			return status.Error(codes.Internal, "internal")
		}

		res = resp
		return nil
	})
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("application %d were created", res.ID)
	err = s.repository.History.NewEvent(ctx, req.UserID, msg)
	if err != nil {
		logger.Error.Println("error writing history message:", msg)
	}
	return res, nil
}
func (s *Application) UserList(ctx context.Context, req *dto.ApplicationUserListRequest) ([]*entity.ApplicationPublic, error) {
	return s.repository.Application.List(ctx, &dto.ApplicationListRequest{
		UserID:   &req.UserID,
		StatusID: req.StatusID,
	})
}
func (s *Application) UserItem(ctx context.Context, req *dto.ApplicationUserItemRequest) (*entity.ApplicationPublic, error) {
	res, err := s.repository.Application.Item(ctx, &dto.ApplicationItemRequest{
		UserID:        &req.UserID,
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, "application not exist")
	}
	return res, nil
}
func (s *Application) DeleteItem(ctx context.Context, req *dto.ApplicationItemDeleteRequest) error {
	err := s.repository.TxManager().ReadWriteTx(ctx, func(ctx context.Context) error {
		item, err := s.repository.Application.Item(ctx, &dto.ApplicationItemRequest{
			UserID:        &req.UserID,
			ApplicationID: req.ApplicationID,
		})
		if err != nil {
			logger.Error.Println("error get application:", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		if item == nil {
			return status.Error(codes.InvalidArgument, "application not found")
		}
		if item.StatusID != int32(pb.ApplicationStatusID_created) {
			return status.Error(codes.InvalidArgument, "application can't be deleted")
		}

		_, err = s.repository.Application.UpdateStatus(ctx, &dto.ApplicationUpdateStatusRequest{
			ApplicationID: item.ID,
			StatusID:      int32(pb.ApplicationStatusID_deleted),
		})
		if err != nil {
			logger.Error.Println("error deleting application status:", err.Error())
			return status.Error(codes.Internal, "internal")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Application) ManagementList(ctx context.Context, req *dto.ApplicationManagementListRequest) ([]*entity.ApplicationPublic, error) {
	return s.repository.Application.List(ctx, &dto.ApplicationListRequest{
		UserID:   req.UserID,
		StatusID: req.StatusID,
	})
}
func (s *Application) ManagementItem(ctx context.Context, req *dto.ApplicationManagementItemRequest) (*entity.ApplicationPublic, error) {
	res, err := s.repository.Application.Item(ctx, &dto.ApplicationItemRequest{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, "application not exist")
	}
	return res, nil
}
func (s *Application) ManagementPrivateItem(ctx context.Context, req *dto.ApplicationManagementItemRequest) (*entity.ApplicationPrivate, error) {
	userID := utils.ContextGetUserID(ctx)

	res, err := s.repository.Application.PrivateItem(ctx, &dto.ApplicationItemRequest{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, "application not exist")
	}

	msg := fmt.Sprintf("get private application_id=%d", req.ApplicationID)
	if err := s.repository.History.NewEvent(ctx, userID, msg); err != nil {
		logger.Error.Println("error writing history message:", msg)
	}
	return res, nil
}
func (s *Application) ManagementUpdateStatus(ctx context.Context, req *dto.ApplicationManagementUpdateStatusRequest) (*entity.ApplicationPublic, error) {
	resp, err := s.repository.Application.UpdateStatus(ctx, &dto.ApplicationUpdateStatusRequest{
		ApplicationID: req.ApplicationID,
		StatusID:      req.StatusID,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (s *Application) ManagementUpdateItem(ctx context.Context, req *dto.ApplicationManagementUpdateItemRequest) (*entity.ApplicationPublic, error) {
	resp, err := s.repository.Application.UpdateItem(ctx, &dto.ApplicationUpdateRequest{
		ApplicationID: req.ApplicationID,
		CurrentMMR:    req.CurrentMMR,
		TargetMMR:     req.TargetMMR,
		Price:         &req.Price,
	})
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("update data application_id=%d", req.ApplicationID)
	if err = s.repository.History.NewEvent(ctx, req.UserID, msg); err != nil {
		logger.Error.Println("error writing history message:", msg)
	}
	return resp, nil
}
func (s *Application) ManagementUpdatePrivate(ctx context.Context, req *dto.ApplicationManagementUpdatePrivateRequest) (*entity.ApplicationPrivate, error) {
	resp, err := s.repository.Application.UpdatePrivate(ctx, &dto.ApplicationUpdatePrivateRequest{
		ApplicationID: req.ApplicationID,
		SteamLogin:    req.SteamLogin,
		SteamPassword: req.SteamPassword,
	})
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("update private application_id=%d", req.ApplicationID)
	if err = s.repository.History.NewEvent(ctx, req.UserID, msg); err != nil {
		logger.Error.Println("error writing history message:", msg)
	}
	return resp, nil
}
