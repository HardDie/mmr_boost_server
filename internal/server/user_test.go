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

func TestUser_Password(t *testing.T) {
	ctx := context.Background()
	serviceUser := mocks.NewIServiceUser(t)
	srv := newUser(service.NewService(nil, nil, nil, serviceUser))

	serviceUser.On("UpdatePassword",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.UserUpdatePasswordRequest{OldPassword: "test", NewPassword: "new"},
		int32(1),
	).
		Return(nil)
	serviceUser.On("UpdatePassword",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.UserUpdatePasswordRequest{OldPassword: "internal", NewPassword: "error"},
		int32(1),
	).
		Return(status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetUserID(ctx, 1)

	tests := []struct {
		name    string
		req     *pb.PasswordRequest
		resp    *emptypb.Empty
		errCode codes.Code
	}{
		{
			"valid",
			&pb.PasswordRequest{OldPassword: "test", NewPassword: "new"},
			&emptypb.Empty{},
			codes.OK,
		},
		{
			"validation error",
			&pb.PasswordRequest{OldPassword: "test", NewPassword: "test"},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.PasswordRequest{OldPassword: "internal", NewPassword: "error"},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.Password(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateEmptyResponse(t, resp, tt.resp)
		})
	}
}

func TestUser_SteamID(t *testing.T) {
	ctx := context.Background()
	serviceUser := mocks.NewIServiceUser(t)
	srv := newUser(service.NewService(nil, nil, nil, serviceUser))

	serviceUser.On("UpdateSteamID",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.UserUpdateSteamIDRequest{SteamID: "666"},
		int32(1),
	).
		Return(&entity.User{
			ID:          1,
			Email:       "test@mail.com",
			Username:    "test",
			RoleID:      1,
			SteamID:     utils.Allocate("666"),
			IsActivated: true,
		}, nil)
	serviceUser.On("UpdateSteamID",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.UserUpdateSteamIDRequest{SteamID: "000"},
		int32(1),
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetUserID(ctx, 1)

	tests := []struct {
		name    string
		req     *pb.SteamIDRequest
		resp    *pb.SteamIDResponse
		errCode codes.Code
	}{
		{
			"valid",
			&pb.SteamIDRequest{SteamId: "666"},
			&pb.SteamIDResponse{Data: &pb.UserObject{
				Id:          1,
				Email:       "test@mail.com",
				Username:    "test",
				RoleId:      1,
				SteamId:     utils.Allocate("666"),
				IsActivated: true,
				CreatedAt:   timestamppb.New(time.Time{}),
				UpdatedAt:   timestamppb.New(time.Time{}),
			}},
			codes.OK,
		},
		{
			"validation error",
			&pb.SteamIDRequest{SteamId: "not numeric"},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.SteamIDRequest{SteamId: "000"},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.SteamID(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}
