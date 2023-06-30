package service

import (
	"context"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
)

type StatusHistory struct {
	repository *postgres.Postgres
}

func NewStatusHistory(repository *postgres.Postgres) *StatusHistory {
	return &StatusHistory{
		repository: repository,
	}
}

func (s *StatusHistory) StatusHistory(ctx context.Context, req *dto.StatusHistoryListRequest) ([]*entity.StatusHistory, error) {
	items, err := s.repository.StatusHistory.List(ctx, req.ApplicationID)
	if err != nil {
		return nil, err
	}
	return items, nil
}
