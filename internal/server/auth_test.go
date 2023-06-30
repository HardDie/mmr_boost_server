package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/mocks"
	"github.com/HardDie/mmr_boost_server/internal/service"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

func TestAuth_Register(t *testing.T) {
	ctx := context.Background()
	serviceAuth := mocks.NewIServiceAuth(t)
	srv := newAuth(service.NewService(nil, serviceAuth, nil, nil, nil, nil))

	serviceAuth.On("Register",
		mock.AnythingOfType("*context.emptyCtx"),
		&dto.AuthRegisterRequest{Username: "test", Password: "test", Email: "test@mail.com"},
	).
		Return(nil)
	serviceAuth.On("Register",
		mock.AnythingOfType("*context.emptyCtx"),
		&dto.AuthRegisterRequest{Username: "internal", Password: "internal", Email: "test@mail.com"},
	).
		Return(status.Error(codes.Internal, "internal"))

	tests := []struct {
		name    string
		req     *pb.RegisterRequest
		resp    *emptypb.Empty
		errCode codes.Code
	}{
		{
			"valid",
			&pb.RegisterRequest{Username: "test", Password: "test", Email: "test@mail.com"},
			&emptypb.Empty{},
			codes.OK,
		},
		{
			"validation error",
			&pb.RegisterRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.RegisterRequest{Username: "internal", Password: "internal", Email: "test@mail.com"},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.Register(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateEmptyResponse(t, resp, tt.resp)
		})
	}
}

func TestAuth_Login(t *testing.T) {
	ctx := context.Background()
	serviceAuth := mocks.NewIServiceAuth(t)
	srv := newAuth(service.NewService(nil, serviceAuth, nil, nil, nil, nil))

	serviceAuth.On("Login",
		mock.AnythingOfType("*context.emptyCtx"),
		&dto.AuthLoginRequest{Username: "test", Password: "test"},
	).
		Return(&entity.User{ID: 1}, nil)
	serviceAuth.On("Login",
		mock.AnythingOfType("*context.emptyCtx"),
		&dto.AuthLoginRequest{Username: "test2", Password: "test2"},
	).
		Return(&entity.User{ID: 2}, nil)
	serviceAuth.On("Login",
		mock.AnythingOfType("*context.emptyCtx"),
		&dto.AuthLoginRequest{Username: "internal", Password: "internal"},
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	serviceAuth.On("GenerateCookie",
		mock.AnythingOfType("*context.emptyCtx"),
		int32(1),
	).
		Return(&entity.AccessToken{TokenHash: "session"}, nil)
	serviceAuth.On("GenerateCookie",
		mock.AnythingOfType("*context.emptyCtx"),
		int32(2),
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	tests := []struct {
		name    string
		req     *pb.LoginRequest
		resp    *emptypb.Empty
		errCode codes.Code
	}{
		{
			"valid",
			&pb.LoginRequest{Username: "test", Password: "test"},
			&emptypb.Empty{},
			codes.OK,
		},
		{
			"validation error",
			&pb.LoginRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"login internal error",
			&pb.LoginRequest{Username: "internal", Password: "internal"},
			nil,
			codes.Internal,
		},
		{
			"cookie internal error",
			&pb.LoginRequest{Username: "test2", Password: "test2"},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.Login(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateEmptyResponse(t, resp, tt.resp)
		})
	}
}

func TestAuth_ValidateEmail(t *testing.T) {
	ctx := context.Background()
	serviceAuth := mocks.NewIServiceAuth(t)
	srv := newAuth(service.NewService(nil, serviceAuth, nil, nil, nil, nil))

	uuid, err := utils.UUIDGenerate()
	if err != nil {
		t.Fatal("error generate uuid:", err.Error())
	}

	serviceAuth.On("ValidateEmail",
		mock.AnythingOfType("*context.emptyCtx"),
		uuid,
	).
		Return(nil)
	serviceAuth.On("ValidateEmail",
		mock.AnythingOfType("*context.emptyCtx"),
		"deadc0de-dead-c0de-dead-c0dedeadc0de",
	).
		Return(status.Error(codes.Internal, "internal"))

	tests := []struct {
		name    string
		req     *pb.ValidateEmailRequest
		resp    *emptypb.Empty
		errCode codes.Code
	}{
		{
			"valid",
			&pb.ValidateEmailRequest{Code: uuid},
			&emptypb.Empty{},
			codes.OK,
		},
		{
			"validation error",
			&pb.ValidateEmailRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.ValidateEmailRequest{Code: "deadc0de-dead-c0de-dead-c0dedeadc0de"},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.ValidateEmail(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateEmptyResponse(t, resp, tt.resp)
		})
	}
}

func TestAuth_SendValidationEmail(t *testing.T) {
	ctx := context.Background()
	serviceAuth := mocks.NewIServiceAuth(t)
	srv := newAuth(service.NewService(nil, serviceAuth, nil, nil, nil, nil))

	serviceAuth.On("SendValidationEmail",
		mock.AnythingOfType("*context.emptyCtx"),
		"test",
	).
		Return(nil)
	serviceAuth.On("SendValidationEmail",
		mock.AnythingOfType("*context.emptyCtx"),
		"internal",
	).
		Return(status.Error(codes.Internal, "internal"))

	tests := []struct {
		name    string
		req     *pb.SendValidationEmailRequest
		resp    *emptypb.Empty
		errCode codes.Code
	}{
		{
			"valid",
			&pb.SendValidationEmailRequest{Username: "test"},
			&emptypb.Empty{},
			codes.OK,
		},
		{
			"validation error",
			&pb.SendValidationEmailRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.SendValidationEmailRequest{Username: "internal"},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.SendValidationEmail(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateEmptyResponse(t, resp, tt.resp)
		})
	}
}

func TestAuth_User(t *testing.T) {
	serviceAuth := mocks.NewIServiceAuth(t)
	srv := newAuth(service.NewService(nil, serviceAuth, nil, nil, nil, nil))

	serviceAuth.On("GetUserInfo",
		mock.AnythingOfType("*context.valueCtx"),
		int32(1),
	).
		Return(&entity.User{
			ID:          1,
			Email:       "test@mail.com",
			Username:    "test",
			RoleID:      1,
			IsActivated: true,
		}, nil)
	serviceAuth.On("GetUserInfo",
		mock.AnythingOfType("*context.valueCtx"),
		int32(2),
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	tests := []struct {
		name    string
		reqCtx  context.Context
		resp    *pb.UserResponse
		errCode codes.Code
	}{
		{
			"valid",
			utils.ContextSetUserID(context.Background(), 1),
			&pb.UserResponse{
				Data: &pb.UserObject{
					Id:          1,
					Email:       "test@mail.com",
					Username:    "test",
					RoleId:      1,
					IsActivated: true,
					CreatedAt:   timestamppb.New(time.Time{}),
					UpdatedAt:   timestamppb.New(time.Time{}),
				},
			},
			codes.OK,
		},
		{
			"internal error",
			utils.ContextSetUserID(context.Background(), 2),
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.User(tt.reqCtx, &emptypb.Empty{})
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}

func TestAuth_Logout(t *testing.T) {
	serviceAuth := mocks.NewIServiceAuth(t)
	srv := newAuth(service.NewService(nil, serviceAuth, nil, nil, nil, nil))

	serviceAuth.On("Logout",
		mock.AnythingOfType("*context.valueCtx"),
		int32(1),
	).
		Return(nil)
	serviceAuth.On("Logout",
		mock.AnythingOfType("*context.valueCtx"),
		int32(2),
	).
		Return(status.Error(codes.Internal, "internal"))

	tests := []struct {
		name    string
		reqCtx  context.Context
		resp    *emptypb.Empty
		errCode codes.Code
	}{
		{
			"valid",
			utils.ContextSetSession(context.Background(), &entity.AccessToken{ID: 1}),
			&emptypb.Empty{},
			codes.OK,
		},
		{
			"internal error",
			utils.ContextSetSession(context.Background(), &entity.AccessToken{ID: 2}),
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.Logout(tt.reqCtx, &emptypb.Empty{})
			validateError(t, err, tt.errCode)
			validateEmptyResponse(t, resp, tt.resp)
		})
	}
}
