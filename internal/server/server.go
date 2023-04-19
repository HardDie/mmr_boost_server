package server

import (
	"context"
	"time"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/middleware"
	"github.com/HardDie/mmr_boost_server/internal/service"
)

type Server struct {
	application
	auth
	system
	user

	authMiddleware    *middleware.AuthMiddleware
	timeoutMiddleware *middleware.TimeoutRequestMiddleware
}

func NewServer(config config.Config, srvc *service.Service) *Server {
	return &Server{
		application: newApplication(srvc),
		auth:        newAuth(srvc),
		system:      newSystem(srvc),
		user:        newUser(srvc),

		authMiddleware:    middleware.NewAuthMiddleware(srvc),
		timeoutMiddleware: middleware.NewTimeoutRequestMiddleware(time.Duration(config.Http.RequestTimeout) * time.Second),
	}
}

func customOutgoingHeaderMatcher(h string) (string, bool) {
	return h, true
}

func (s *Server) Register(router *mux.Router) {
	ctx := context.Background()
	router.Use(
		middleware.LoggerMiddleware,
		s.authMiddleware.RequestMiddleware,
	)

	// Create grpc mux
	grpcMux := runtime.NewServeMux(
		// Fix "Grpc-Metadata-" prefix for outgoing headers
		runtime.WithOutgoingHeaderMatcher(customOutgoingHeaderMatcher),
	)

	systemRouter := router.PathPrefix("/system").Subrouter()
	s.system.RegisterPublicRouter(systemRouter, middleware.CorsMiddleware, s.timeoutMiddleware.RequestMiddleware)

	err := s.application.RegisterHTTP(ctx, grpcMux)
	if err != nil {
		logger.Error.Fatal("error register application", err.Error())
	}

	err = s.auth.RegisterHTTP(ctx, grpcMux)
	if err != nil {
		logger.Error.Fatal("error register auth", err.Error())
	}

	err = s.user.RegisterHTTP(ctx, grpcMux)
	if err != nil {
		logger.Error.Fatal("error register user", err.Error())
	}

	router.PathPrefix("/").Handler(grpcMux)
}
