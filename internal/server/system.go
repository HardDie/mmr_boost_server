package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/HardDie/mmr_boost_server/internal/service"
	pb "github.com/HardDie/mmr_boost_server/pkg/proto/server"
)

type system struct {
	service *service.Service
	pb.UnimplementedSystemServer
}

func newSystem(service *service.Service) system {
	return system{
		service: service,
	}
}

func (s *system) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	return pb.RegisterSystemHandlerServer(ctx, mux, s)
}

func (s *system) Swagger(ctx context.Context, _ *emptypb.Empty) (*httpbody.HttpBody, error) {
	data, err := s.service.SystemGetSwagger()
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &httpbody.HttpBody{
		ContentType: "application/yaml",
		Data:        data,
	}, nil
}
