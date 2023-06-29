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

var calibrationPrice float64 = 2000
var price = []struct {
	minMMR int32
	maxMMR int32
	price  float64
}{
	{0, 1500, 160. / 100},
	{1500, 2500, 184. / 100},
	{2500, 3500, 240. / 100},
	{3500, 4000, 296. / 100},
	{4000, 4500, 336. / 100},
	{4500, 4800, 352. / 100},
	{4800, 4900, 368. / 100},
	{4900, 5000, 384. / 100},
	{5000, 5100, 392. / 100},
	{5100, 5200, 420. / 100},
	{5200, 5300, 460. / 100},
	{5300, 5400, 480. / 100},
	{5400, 5500, 520. / 100},
	{5500, 5600, 560. / 100},
	{5600, 5700, 600. / 100},
	{5700, 5800, 640. / 100},
	{5800, 5900, 680. / 100},
	{5900, 6000, 720. / 100},
	{6000, 6100, 800. / 100},
	{6100, 6200, 860. / 100},
	{6200, 6300, 920. / 100},
	{6300, 6400, 960. / 100},
	{6400, 6500, 1000. / 100},
	{6500, 6800, 1040. / 100},
	{6800, 6900, 1080. / 100},
	{6900, 7000, 1200. / 100},
	{7000, 7500, 1400. / 100},
	{7500, 8000, 1500. / 100},
}

func (s *Price) Price(_ context.Context, req *dto.PriceRequest) (float64, error) {
	if req.TypeID == int32(pb.ApplicationTypeID_calibration) {
		return calibrationPrice, nil
	}
	if req.TypeID != int32(pb.ApplicationTypeID_boost_mmr) {
		return 0, status.Error(codes.InvalidArgument, "unknown service type")
	}
	if req.TargetMmr <= req.CurrentMmr {
		return 0, status.Error(codes.InvalidArgument, "invalid mmr values")
	}

	var res float64

	for _, p := range price {
		if p.maxMMR < req.CurrentMmr {
			continue
		}
		if p.minMMR > req.TargetMmr {
			break
		}

		minMMR := p.minMMR
		if req.CurrentMmr >= p.minMMR &&
			req.CurrentMmr < p.maxMMR {
			minMMR = req.CurrentMmr
		}
		maxMMR := p.maxMMR
		if req.TargetMmr >= p.minMMR &&
			req.TargetMmr < p.maxMMR {
			maxMMR = req.TargetMmr
		}

		res += float64(maxMMR-minMMR) * p.price
	}

	return res, nil
}
