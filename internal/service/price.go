package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type Price struct {
	repository *postgres.Postgres
}

func NewPrice(repository *postgres.Postgres) *Price {
	return &Price{
		repository: repository,
	}
}

func (s *Price) Price(ctx context.Context, req *dto.PriceRequest) (float64, error) {
	switch req.TypeID {
	case int32(pb.ApplicationTypeID_boost_mmr):
		val := float64(req.TargetMmr-req.CurrentMmr) / 100 * 300
		return val, nil
	case int32(pb.ApplicationTypeID_calibration):
		return 1000, nil
	}
	return 0, status.Error(codes.InvalidArgument, "unknown service type")
}
