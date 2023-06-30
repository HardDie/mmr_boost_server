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

type statusHistory struct {
	service *service.Service
	pb.UnimplementedStatusServer
}

func newStatusHistory(service *service.Service) statusHistory {
	return statusHistory{
		service: service,
	}
}

func (s *statusHistory) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	return pb.RegisterStatusHandlerServer(ctx, mux, s)
}

func (s *statusHistory) StatusHistory(ctx context.Context, req *pb.StatusHistoryRequest) (*pb.StatusHistoryResponse, error) {
	r := &dto.StatusHistoryListRequest{
		ApplicationID: req.Id,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.StatusHistory.StatusHistory(ctx, r)
	if err != nil {
		return nil, err
	}

	var data []*pb.StatusHistory
	for _, item := range resp {
		data = append(data, StatusHistoryToPb(item))
	}
	return &pb.StatusHistoryResponse{
		Data: data,
	}, nil
}
