package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/repository/encrypt"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type Application struct {
	repository        *postgres.Postgres
	encryptRepository *encrypt.Encrypt
}

func NewApplication(repository *postgres.Postgres, encrypt *encrypt.Encrypt) *Application {
	return &Application{
		repository:        repository,
		encryptRepository: encrypt,
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
	err = s.repository.StatusHistory.NewEvent(ctx, &dto.StatusHistoryNewEventRequest{
		UserID:        req.UserID,
		ApplicationID: req.ApplicationID,
		NewStatusID:   int32(pb.ApplicationStatusID_deleted),
	})
	if err != nil {
		logger.Error.Printf("error writing status_history message: user_id=%d, application_id=%d new_status_id=%d",
			req.UserID, req.ApplicationID, int32(pb.ApplicationStatusID_deleted))
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
func (s *Application) ManagementPrivateItem(ctx context.Context, req *dto.ApplicationManagementPrivateItemRequest) (*entity.ApplicationPrivate, error) {
	res, err := s.repository.Application.PrivateItem(ctx, &dto.ApplicationItemRequest{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, "application not exist")
	}

	if res.SteamLogin != nil && len(*res.SteamLogin) > 0 {
		login, err := s.encryptRepository.Decrypt(*res.SteamLogin)
		if err != nil {
			logger.Error.Println("error decrypt steam login:", err.Error())
			return nil, status.Error(codes.Internal, "internal")
		}
		*res.SteamLogin = login
	}
	if res.SteamPassword != nil && len(*res.SteamPassword) > 0 {
		password, err := s.encryptRepository.Decrypt(*res.SteamPassword)
		if err != nil {
			logger.Error.Println("error decrypt steam password:", err.Error())
			return nil, status.Error(codes.Internal, "internal")
		}
		*res.SteamPassword = password
	}

	msg := fmt.Sprintf("get private application_id=%d", req.ApplicationID)
	if err := s.repository.History.NewEvent(ctx, req.UserID, msg); err != nil {
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
	err = s.repository.StatusHistory.NewEvent(ctx, &dto.StatusHistoryNewEventRequest{
		UserID:        req.UserID,
		ApplicationID: req.ApplicationID,
		NewStatusID:   req.StatusID,
	})
	if err != nil {
		logger.Error.Printf("error writing status_history message: user_id=%d, application_id=%d new_status_id=%d",
			req.UserID, req.ApplicationID, req.StatusID)
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
	login, err := s.encryptRepository.Encrypt(req.SteamLogin)
	if err != nil {
		logger.Error.Println("error encrypt steam login:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	password, err := s.encryptRepository.Encrypt(req.SteamPassword)
	if err != nil {
		logger.Error.Println("error encrypt steam password:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}

	resp, err := s.repository.Application.UpdatePrivate(ctx, &dto.ApplicationUpdatePrivateRequest{
		ApplicationID: req.ApplicationID,
		SteamLogin:    login,
		SteamPassword: password,
	})
	if err != nil {
		return nil, err
	}

	*resp.SteamLogin = req.SteamLogin
	*resp.SteamPassword = req.SteamPassword

	msg := fmt.Sprintf("update private application_id=%d", req.ApplicationID)
	if err = s.repository.History.NewEvent(ctx, req.UserID, msg); err != nil {
		logger.Error.Println("error writing history message:", msg)
	}
	return resp, nil
}
