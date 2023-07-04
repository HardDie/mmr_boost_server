package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/service"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type user struct {
	service *service.Service
	pb.UnimplementedUserServer
}

func newUser(service *service.Service) user {
	return user{
		service: service,
	}
}

func (s *user) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	return pb.RegisterUserHandlerServer(ctx, mux, s)
}

func (s *user) Password(ctx context.Context, req *pb.PasswordRequest) (*emptypb.Empty, error) {
	userID, err := utils.ContextGetUserID(ctx)
	if err != nil {
		logger.Error.Printf("userID not found in context")
		return nil, err
	}

	r := &dto.UserUpdatePasswordRequest{
		NewPassword: req.NewPassword,
		OldPassword: req.OldPassword,
	}
	err = getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.User.UpdatePassword(ctx, r, userID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (s *user) SteamID(ctx context.Context, req *pb.SteamIDRequest) (*pb.SteamIDResponse, error) {
	userID, err := utils.ContextGetUserID(ctx)
	if err != nil {
		logger.Error.Printf("userID not found in context")
		return nil, err
	}

	r := &dto.UserUpdateSteamIDRequest{
		SteamID: req.SteamId,
	}
	err = getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	u, err := s.service.User.UpdateSteamID(ctx, r, userID)
	if err != nil {
		return nil, err
	}

	return &pb.SteamIDResponse{
		Data: UserToPb(u),
	}, nil
}
