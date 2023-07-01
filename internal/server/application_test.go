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
	"github.com/HardDie/mmr_boost_server/internal/utils"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

func TestApplication_Create(t *testing.T) {
	ctx := context.Background()
	m, s := initServerObject(t)
	srv := newApplication(s)

	m.price.On("Price",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.PriceRequest{
			TypeID:     1,
			CurrentMmr: 1000,
			TargetMmr:  2000,
		},
	).Return(int32(150), nil)
	m.price.On("Price",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.PriceRequest{
			TypeID:     1,
			CurrentMmr: 0,
			TargetMmr:  10,
		},
	).Return(int32(0), status.Error(codes.Internal, "internal"))

	m.application.On("Create",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationCreateRequest{
			UserID:     1,
			TypeID:     1,
			CurrentMMR: 1000,
			TargetMMR:  2000,
			TgContact:  "testuser",
			Price:      150,
		},
	).Return(&entity.ApplicationPublic{
		ID:         1,
		UserID:     1,
		StatusID:   1,
		TypeID:     1,
		CurrentMMR: 1000,
		TargetMMR:  2000,
		TgContact:  "testuser",
		Price:      150,
	}, nil)
	m.application.On("Create",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationCreateRequest{
			UserID:     1,
			TypeID:     1,
			CurrentMMR: 1000,
			TargetMMR:  2000,
			TgContact:  "internal",
			Price:      150,
		},
	).Return(nil, status.Error(codes.Internal, "internal"))

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
					Price:      150,
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
			"price internal",
			&pb.CreateRequest{TypeId: 1, CurrentMmr: 0, TargetMmr: 10, TgContact: "price"},
			nil,
			codes.Internal,
		},
		{
			"internal error",
			&pb.CreateRequest{TypeId: 1, CurrentMmr: 1000, TargetMmr: 2000, TgContact: "internal"},
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

	m, s := initServerObject(t)
	srv := newApplication(s)

	m.application.On("UserList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserListRequest{UserID: 1},
	).Return([]*entity.ApplicationPublic{
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
	m.application.On("UserList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserListRequest{UserID: 1, StatusID: utils.Allocate[int32](3)},
	).Return(nil, status.Error(codes.Internal, "internal"))

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
			&pb.GetListRequest{StatusId: utils.Allocate[pb.ApplicationStatusID](-1)},
			nil,
			codes.InvalidArgument,
		},
		{
			"internal error",
			&pb.GetListRequest{StatusId: utils.Allocate[pb.ApplicationStatusID](3)},
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
	m, s := initServerObject(t)
	srv := newApplication(s)

	m.application.On("UserItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserItemRequest{UserID: 1, ApplicationID: 1},
	).Return(&entity.ApplicationPublic{
		ID:         1,
		UserID:     1,
		StatusID:   1,
		TypeID:     1,
		CurrentMMR: 1000,
		TargetMMR:  2000,
		TgContact:  "testuser",
	}, nil)
	m.application.On("UserItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationUserItemRequest{UserID: 1, ApplicationID: 2},
	).Return(nil, status.Error(codes.Internal, "internal"))

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
	m, s := initServerObject(t)
	srv := newApplication(s)

	m.application.On("ManagementList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementListRequest{},
	).Return([]*entity.ApplicationPublic{
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
	m.application.On("ManagementList",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementListRequest{UserID: utils.Allocate[int32](3)},
	).Return(nil, status.Error(codes.Internal, "internal"))

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
			&pb.GetManagementListRequest{StatusId: utils.Allocate[pb.ApplicationStatusID](-1)},
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
	m, s := initServerObject(t)
	srv := newApplication(s)

	m.application.On("ManagementItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementItemRequest{ApplicationID: 1},
	).Return(&entity.ApplicationPublic{
		ID:         1,
		UserID:     1,
		StatusID:   1,
		TypeID:     1,
		CurrentMMR: 1000,
		TargetMMR:  2000,
		TgContact:  "testuser",
	}, nil)
	m.application.On("ManagementItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementItemRequest{ApplicationID: 3},
	).Return(nil, status.Error(codes.Internal, "internal"))

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
	m, s := initServerObject(t)
	srv := newApplication(s)

	m.application.On("ManagementPrivateItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementPrivateItemRequest{ApplicationID: 1, UserID: 1},
	).Return(&entity.ApplicationPrivate{
		ID:            1,
		SteamLogin:    utils.Allocate("testlogin"),
		SteamPassword: utils.Allocate("testpassword"),
	}, nil)
	m.application.On("ManagementPrivateItem",
		mock.AnythingOfType("*context.valueCtx"),
		&dto.ApplicationManagementPrivateItemRequest{ApplicationID: 3, UserID: 1},
	).Return(nil, status.Error(codes.Internal, "internal"))

	// Set user for context
	ctx = utils.ContextSetUserID(ctx, 1)
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
