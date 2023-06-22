package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/mocks"
	"github.com/HardDie/mmr_boost_server/internal/service"
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

func TestApplication_Create(t *testing.T) {
	ctx := context.Background()
	serviceApplication := mocks.NewIServiceApplication(t)
	srv := newApplication(service.NewService(serviceApplication, nil, nil, nil))

	serviceApplication.On("Create",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationCreateRequest{
			UserID:     1,
			TypeID:     1,
			CurrentMMR: 1000,
			TargetMMR:  2000,
			TgContact:  "testuser",
		},
	).
		Return(&entity.ApplicationPublic{
			ID:         1,
			UserID:     1,
			StatusID:   1,
			TypeID:     1,
			CurrentMMR: 1000,
			TargetMMR:  2000,
			TgContact:  "testuser",
		}, nil)
	serviceApplication.On("Create",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationCreateRequest{
			UserID:     1,
			TypeID:     1,
			CurrentMMR: 10,
			TargetMMR:  666,
			TgContact:  "internal",
		},
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetUserID(ctx, 1)

	tests := []struct {
		name    string
		req     *pb.CreateRequest
		resp    *pb.CreateResponse
		errCode codes.Code
	}{
		{
			"valid",
			&pb.CreateRequest{
				TypeId:     1,
				CurrentMmr: 1000,
				TargetMmr:  2000,
				TgContact:  "testuser",
			},
			&pb.CreateResponse{
				Data: &pb.PublicApplicationObject{
					Id:         1,
					UserId:     1,
					StatusId:   1,
					TypeId:     1,
					CurrentMmr: 1000,
					TargetMmr:  2000,
					TgContact:  "testuser",
					CreatedAt:  timestamppb.New(time.Time{}),
					UpdatedAt:  timestamppb.New(time.Time{}),
				},
			},
			codes.OK,
		},
		{
			"validation error",
			&pb.CreateRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.CreateRequest{TypeId: 1, CurrentMmr: 10, TargetMmr: 666, TgContact: "internal"},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.Create(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}

func TestApplication_GetList(t *testing.T) {
	ctx := context.Background()
	serviceApplication := mocks.NewIServiceApplication(t)
	srv := newApplication(service.NewService(serviceApplication, nil, nil, nil))

	serviceApplication.On("UserList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserListRequest{UserID: 1},
	).
		Return([]*entity.ApplicationPublic{
			{
				ID:         1,
				UserID:     1,
				StatusID:   1,
				TypeID:     1,
				CurrentMMR: 1000,
				TargetMMR:  2000,
				TgContact:  "testuser",
			},
		}, nil)
	serviceApplication.On("UserList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserListRequest{UserID: 1, StatusID: utils.Allocate[int32](3)},
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetUserID(ctx, 1)

	tests := []struct {
		name    string
		req     *pb.GetListRequest
		resp    *pb.GetListResponse
		errCode codes.Code
	}{
		{
			"valid",
			&pb.GetListRequest{},
			&pb.GetListResponse{
				Data: []*pb.PublicApplicationObject{
					{
						Id:         1,
						UserId:     1,
						StatusId:   1,
						TypeId:     1,
						CurrentMmr: 1000,
						TargetMmr:  2000,
						TgContact:  "testuser",
						CreatedAt:  timestamppb.New(time.Time{}),
						UpdatedAt:  timestamppb.New(time.Time{}),
					},
				},
			},
			codes.OK,
		},
		{
			"validation error",
			&pb.GetListRequest{StatusId: utils.Allocate[int32](-1)},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.GetListRequest{StatusId: utils.Allocate[int32](3)},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.GetList(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}

func TestApplication_GetItem(t *testing.T) {
	ctx := context.Background()
	serviceApplication := mocks.NewIServiceApplication(t)
	srv := newApplication(service.NewService(serviceApplication, nil, nil, nil))

	serviceApplication.On("UserItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserItemRequest{UserID: 1, ApplicationID: 1},
	).
		Return(&entity.ApplicationPublic{
			ID:         1,
			UserID:     1,
			StatusID:   1,
			TypeID:     1,
			CurrentMMR: 1000,
			TargetMMR:  2000,
			TgContact:  "testuser",
		}, nil)
	serviceApplication.On("UserItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserItemRequest{UserID: 1, ApplicationID: 2},
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetUserID(ctx, 1)

	tests := []struct {
		name    string
		req     *pb.GetItemRequest
		resp    *pb.GetItemResponse
		errCode codes.Code
	}{
		{
			"valid",
			&pb.GetItemRequest{Id: 1},
			&pb.GetItemResponse{
				Data: &pb.PublicApplicationObject{
					Id:         1,
					UserId:     1,
					StatusId:   1,
					TypeId:     1,
					CurrentMmr: 1000,
					TargetMmr:  2000,
					TgContact:  "testuser",
					CreatedAt:  timestamppb.New(time.Time{}),
					UpdatedAt:  timestamppb.New(time.Time{}),
				},
			},
			codes.OK,
		},
		{
			"validation error",
			&pb.GetItemRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.GetItemRequest{Id: 2},
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.GetItem(ctx, tt.req)
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}

func TestApplication_GetManagementList(t *testing.T) {
	ctx := context.Background()
	serviceApplication := mocks.NewIServiceApplication(t)
	srv := newApplication(service.NewService(serviceApplication, nil, nil, nil))

	serviceApplication.On("ManagementList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementListRequest{},
	).
		Return([]*entity.ApplicationPublic{
			{
				ID:         1,
				UserID:     1,
				StatusID:   1,
				TypeID:     1,
				CurrentMMR: 1000,
				TargetMMR:  2000,
				TgContact:  "testuser",
			},
		}, nil)
	serviceApplication.On("ManagementList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementListRequest{UserID: utils.Allocate[int32](3)},
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetUserID(ctx, 1)
	ctx = utils.ContextSetRoleID(ctx, int32(pb.UserRoleID_admin))

	tests := []struct {
		name       string
		req        *pb.GetManagementListRequest
		reqContext context.Context
		resp       *pb.GetManagementListResponse
		errCode    codes.Code
	}{
		{
			"valid",
			&pb.GetManagementListRequest{},
			ctx,
			&pb.GetManagementListResponse{
				Data: []*pb.PublicApplicationObject{
					{
						Id:         1,
						UserId:     1,
						StatusId:   1,
						TypeId:     1,
						CurrentMmr: 1000,
						TargetMmr:  2000,
						TgContact:  "testuser",
						CreatedAt:  timestamppb.New(time.Time{}),
						UpdatedAt:  timestamppb.New(time.Time{}),
					},
				},
			},
			codes.OK,
		},
		{
			"validation error",
			&pb.GetManagementListRequest{StatusId: utils.Allocate[int32](-1)},
			ctx,
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.GetManagementListRequest{UserId: utils.Allocate[int32](3)},
			ctx,
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.GetManagementList(tt.reqContext, tt.req)
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}

func TestApplication_GetManagementItem(t *testing.T) {
	ctx := context.Background()
	serviceApplication := mocks.NewIServiceApplication(t)
	srv := newApplication(service.NewService(serviceApplication, nil, nil, nil))

	serviceApplication.On("ManagementItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementItemRequest{ApplicationID: 1},
	).
		Return(&entity.ApplicationPublic{
			ID:         1,
			UserID:     1,
			StatusID:   1,
			TypeID:     1,
			CurrentMMR: 1000,
			TargetMMR:  2000,
			TgContact:  "testuser",
		}, nil)
	serviceApplication.On("ManagementItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementItemRequest{ApplicationID: 3},
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetRoleID(ctx, int32(pb.UserRoleID_admin))

	tests := []struct {
		name       string
		req        *pb.GetManagementItemRequest
		reqContext context.Context
		resp       *pb.GetManagementItemResponse
		errCode    codes.Code
	}{
		{
			"valid",
			&pb.GetManagementItemRequest{Id: 1},
			ctx,
			&pb.GetManagementItemResponse{
				Data: &pb.PublicApplicationObject{
					Id:         1,
					UserId:     1,
					StatusId:   1,
					TypeId:     1,
					CurrentMmr: 1000,
					TargetMmr:  2000,
					TgContact:  "testuser",
					CreatedAt:  timestamppb.New(time.Time{}),
					UpdatedAt:  timestamppb.New(time.Time{}),
				},
			},
			codes.OK,
		},
		{
			"validation error",
			&pb.GetManagementItemRequest{Id: -1},
			ctx,
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.GetManagementItemRequest{Id: 3},
			ctx,
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.GetManagementItem(tt.reqContext, tt.req)
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}

func TestApplication_GetManagementPrivateItem(t *testing.T) {
	ctx := context.Background()
	serviceApplication := mocks.NewIServiceApplication(t)
	srv := newApplication(service.NewService(serviceApplication, nil, nil, nil))

	serviceApplication.On("ManagementPrivateItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementItemRequest{ApplicationID: 1},
	).
		Return(&entity.ApplicationPrivate{
			ID:            1,
			SteamLogin:    utils.Allocate("testlogin"),
			SteamPassword: utils.Allocate("testpassword"),
		}, nil)
	serviceApplication.On("ManagementPrivateItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementItemRequest{ApplicationID: 3},
	).
		Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetRoleID(ctx, int32(pb.UserRoleID_admin))

	tests := []struct {
		name       string
		req        *pb.GetManagementItemRequest
		reqContext context.Context
		resp       *pb.GetManagementPrivateItemResponse
		errCode    codes.Code
	}{
		{
			"valid",
			&pb.GetManagementItemRequest{Id: 1},
			ctx,
			&pb.GetManagementPrivateItemResponse{
				Data: &pb.PrivateApplicationObject{
					Id:            1,
					SteamLogin:    utils.Allocate("testlogin"),
					SteamPassword: utils.Allocate("testpassword"),
					CreatedAt:     timestamppb.New(time.Time{}),
					UpdatedAt:     timestamppb.New(time.Time{}),
				},
			},
			codes.OK,
		},
		{
			"validation error",
			&pb.GetManagementItemRequest{Id: -1},
			ctx,
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.GetManagementItemRequest{Id: 3},
			ctx,
			nil,
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := srv.GetManagementPrivateItem(tt.reqContext, tt.req)
			validateError(t, err, tt.errCode)
			validateResponse(t, resp, tt.resp)
		})
	}
}
