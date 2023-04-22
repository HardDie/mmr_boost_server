package service

import (
	"context"
	"fmt"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/errs"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
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
	resp, err := s.repository.ApplicationCreate(ctx, req)
	if err != nil {
		logger.Error.Println("error creating new application:", err.Error())
		return nil, errs.InternalError
	}

	msg := fmt.Sprintf("application %d were created", resp.ID)
	err = s.repository.HistoryNewEvent(ctx, req.UserID, msg)
	if err != nil {
		logger.Error.Println("error writing history message:", msg)
	}

	return resp, nil
}

func (s *application) ApplicationUserList(ctx context.Context, req *dto.ApplicationUserListRequest) ([]*entity.ApplicationPublic, error) {
	return s.repository.ApplicationUserList(ctx, &dto.ApplicationListRequest{
		UserID:   &req.UserID,
		StatusID: req.StatusID,
	})
}

func (s *application) ApplicationManagementUserList(ctx context.Context, req *dto.ApplicationManagementUserListRequest) ([]*entity.ApplicationPublic, error) {
	return s.repository.ApplicationUserList(ctx, &dto.ApplicationListRequest{
		UserID:   req.UserID,
		StatusID: req.StatusID,
	})
}
