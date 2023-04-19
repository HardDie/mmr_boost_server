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

type auth struct {
	service *service.Service
	pb.UnimplementedAuthServer
}

func newAuth(service *service.Service) auth {
	return auth{
		service: service,
	}
}

func (s *auth) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	return pb.RegisterAuthHandlerServer(ctx, mux, s)
}

func (s *auth) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	r := &dto.AuthRegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.AuthRegister(ctx, r)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (s *auth) Login(ctx context.Context, req *pb.LoginRequest) (*emptypb.Empty, error) {
	r := &dto.AuthLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	u, err := s.service.AuthLogin(ctx, r)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.service.AuthGenerateCookie(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	utils.SetGRPCSessionCookie(ctx, accessToken.TokenHash)
	return &emptypb.Empty{}, nil
}
func (s *auth) ValidateEmail(ctx context.Context, req *pb.ValidateEmailRequest) (*emptypb.Empty, error) {
	r := &dto.AuthValidateEmailRequest{
		Code: req.Code,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.AuthValidateEmail(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *auth) SendValidationEmail(ctx context.Context, req *pb.SendValidationEmailRequest) (*emptypb.Empty, error) {
	r := &dto.AuthSendValidationEmailRequest{
		Username: req.Username,
	}
	err := getValidator().Struct(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.AuthSendValidationEmail(ctx, r.Username)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *auth) User(ctx context.Context, _ *emptypb.Empty) (*pb.UserResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)

	u, err := s.service.AuthGetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Data: UserToPb(u),
	}, nil
}
func (s *auth) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	session := utils.GetAccessTokenFromContext(ctx)

	err := s.service.AuthLogout(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	utils.DeleteGRPCSessionCookie(ctx)
	return &emptypb.Empty{}, nil
}
