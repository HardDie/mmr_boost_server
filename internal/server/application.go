package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

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
	userID := utils.ContextGetUserID(ctx)

	r := &dto.ApplicationCreateRequest{
		UserID:     userID,
		TypeID:     int32(req.TypeId),
		CurrentMMR: req.CurrentMmr,
		TargetMMR:  req.TargetMmr,
		TgContact:  req.TgContact,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Calculate price for application
	r.Price, err = s.service.Price.Price(ctx, &dto.PriceRequest{
		TypeID:     r.TypeID,
		CurrentMmr: r.CurrentMMR,
		TargetMmr:  r.TargetMMR,
	})
	if err != nil {
		return nil, err
	}

	resp, err := s.service.Application.Create(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Data: ApplicationPublicToPb(resp),
	}, nil
}
func (s *application) GetList(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	userID := utils.ContextGetUserID(ctx)

	r := &dto.ApplicationUserListRequest{
		UserID:   userID,
		StatusID: utils.ToInt32(req.StatusId),
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.Application.UserList(ctx, r)
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
func (s *application) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	userID := utils.ContextGetUserID(ctx)

	r := &dto.ApplicationUserItemRequest{
		UserID:        userID,
		ApplicationID: req.Id,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.Application.UserItem(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.GetItemResponse{
		Data: ApplicationPublicToPb(resp),
	}, nil
}
func (s *application) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*emptypb.Empty, error) {
	userID := utils.ContextGetUserID(ctx)

	r := &dto.ApplicationItemDeleteRequest{
		ApplicationID: req.Id,
		UserID:        userID,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.Application.DeleteItem(ctx, r)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *application) GetManagementList(ctx context.Context, req *pb.GetManagementListRequest) (*pb.GetManagementListResponse, error) {
	r := &dto.ApplicationManagementListRequest{
		UserID:   req.UserId,
		StatusID: utils.ToInt32(req.StatusId),
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.Application.ManagementList(ctx, r)
	if err != nil {
		return nil, err
	}

	var data []*pb.PublicApplicationObject
	for _, item := range resp {
		data = append(data, ApplicationPublicToPb(item))
	}
	return &pb.GetManagementListResponse{
		Data: data,
	}, nil
}
func (s *application) GetManagementItem(ctx context.Context, req *pb.GetManagementItemRequest) (*pb.GetManagementItemResponse, error) {
	r := &dto.ApplicationManagementItemRequest{
		ApplicationID: req.Id,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.Application.ManagementItem(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.GetManagementItemResponse{
		Data: ApplicationPublicToPb(resp),
	}, nil
}
func (s *application) GetManagementPrivateItem(ctx context.Context, req *pb.GetManagementItemRequest) (*pb.GetManagementPrivateItemResponse, error) {
	r := &dto.ApplicationManagementItemRequest{
		ApplicationID: req.Id,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.Application.ManagementPrivateItem(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.GetManagementPrivateItemResponse{
		Data: ApplicationPrivateToPb(resp),
	}, nil
}
func (s *application) UpdateManagementItemStatus(ctx context.Context, req *pb.UpdateManagementItemStatusRequest) (*pb.UpdateManagementItemStatusResponse, error) {
	r := &dto.ApplicationManagementUpdateStatusRequest{
		ApplicationID: req.Id,
		StatusID:      int32(req.StatusId),
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.service.Application.ManagementUpdateStatus(ctx, r)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateManagementItemStatusResponse{
		Data: ApplicationPublicToPb(resp),
	}, nil
}
