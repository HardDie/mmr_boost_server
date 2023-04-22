package service

import (
	"context"
	"fmt"

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
	userID := utils.ContextGetUserID(ctx)

	resp, err := s.repository.ApplicationCreate(ctx, req, userID)
	if err != nil {
		logger.Error.Println("error creating new application:", err.Error())
		return nil, errs.InternalError
	}

	msg := fmt.Sprintf("application %d were created", resp.ID)
	err = s.repository.HistoryNewEvent(ctx, userID, msg)
	if err != nil {
		logger.Error.Println("error writing history message:", msg)
	}

	return resp, nil
}

func (s *application) ApplicationUserList(ctx context.Context, req *dto.ApplicationUserListRequest) ([]*entity.ApplicationPublic, error) {
	return s.repository.ApplicationUserList(ctx, req)
}
