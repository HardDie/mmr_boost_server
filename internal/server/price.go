package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/service"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type price struct {
	service *service.Service
	pb.UnimplementedPriceServer
}

func newPrice(service *service.Service) price {
	return price{
		service: service,
	}
}

func (s *price) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	return pb.RegisterPriceHandlerServer(ctx, mux, s)
}

func (s *price) Price(ctx context.Context, req *pb.PriceRequest) (*pb.PriceResponse, error) {
	r := &dto.PriceRequest{
		TypeID:     int32(req.TypeId),
		CurrentMmr: req.GetCurrentMmr(),
		TargetMmr:  req.GetTargetMmr(),
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.service.Price.Price(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.PriceResponse{
		Price: res,
	}, nil
}
