package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/service"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type application struct {
	service *service.Service
	pb.UnimplementedApplicationServer
}

func newApplication(service *service.Service) application {
	return application{
		service: service,
	}
}

func (s *application) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	return pb.RegisterApplicationHandlerServer(ctx, mux, s)
}

func (s *application) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	r := &dto.ApplicationCreateRequest{
		TypeID:     req.TypeId,
		CurrentMMR: req.CurrentMmr,
		TargetMMR:  req.TargetMmr,
		TgContact:  req.TgContact,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.ApplicationCreate(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Data: ApplicationPublicToPb(resp),
	}, nil
}

func (s *application) GetList(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)

	r := &dto.ApplicationUserListRequest{
		UserID:   userID,
		StatusID: req.StatusId,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.ApplicationUserList(ctx, r)
	if err != nil {
		return nil, err
	}

	var data []*pb.PublicApplicationObject
	for _, item := range resp {
		data = append(data, ApplicationPublicToPb(item))
	}
	return &pb.GetListResponse{
		Data: data,
	}, nil
}
