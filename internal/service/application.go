package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/errs"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/utils"
)

type application struct {
	repository *postgres.Postgres
}

func newApplication(repository *postgres.Postgres) application {
	return application{
		repository: repository,
	}
}

func (s *application) ApplicationCreate(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error) {
	resp, err := s.repository.Application.Create(ctx, req)
	if err != nil {
		logger.Error.Println("error creating new application:", err.Error())
		return nil, errs.InternalError
	}

	msg := fmt.Sprintf("application %d were created", resp.ID)
	err = s.repository.History.NewEvent(ctx, req.UserID, msg)
	if err != nil {
		logger.Error.Println("error writing history message:", msg)
	}

	return resp, nil
}
func (s *application) ApplicationUserList(ctx context.Context, req *dto.ApplicationUserListRequest) ([]*entity.ApplicationPublic, error) {
	return s.repository.Application.List(ctx, &dto.ApplicationListRequest{
		UserID:   &req.UserID,
		StatusID: req.StatusID,
	})
}
func (s *application) ApplicationManagementUserList(ctx context.Context, req *dto.ApplicationManagementListRequest) ([]*entity.ApplicationPublic, error) {
	return s.repository.Application.List(ctx, &dto.ApplicationListRequest{
		UserID:   req.UserID,
		StatusID: req.StatusID,
	})
}
func (s *application) ApplicationUserItem(ctx context.Context, req *dto.ApplicationUserItemRequest) (*entity.ApplicationPublic, error) {
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
func (s *application) ApplicationManagementItem(ctx context.Context, req *dto.ApplicationManagementItemRequest) (*entity.ApplicationPublic, error) {
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
func (s *application) ApplicationManagementPrivateItem(ctx context.Context, req *dto.ApplicationManagementItemRequest) (*entity.ApplicationPrivate, error) {
	userID := utils.ContextGetUserID(ctx)

	res, err := s.repository.Application.PrivateItem(ctx, &dto.ApplicationItemRequest{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		if err := s.repository.History.NewEvent(ctx, userID, fmt.Sprintf("error get private application_id=%d", req.ApplicationID)); err != nil {
			logger.Error.Printf("error writing history message: error get private application_id=%d", req.ApplicationID)
		}
		return nil, err
	}

	if res == nil {
		if err := s.repository.History.NewEvent(ctx, userID, fmt.Sprintf("get not exist private application_id=%d", req.ApplicationID)); err != nil {
			logger.Error.Printf("error writing history message: get not exist private application_id=%d", req.ApplicationID)
		}
		return nil, status.Error(codes.NotFound, "application not exist")
	}

	if err := s.repository.History.NewEvent(ctx, userID, fmt.Sprintf("get private application_id=%d", req.ApplicationID)); err != nil {
		logger.Error.Printf("error writing history message: get private application_id=%d", req.ApplicationID)
	}
	return res, nil
}
