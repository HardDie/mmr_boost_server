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
	price

	authMiddleware    *middleware.AuthMiddleware
	timeoutMiddleware *middleware.TimeoutRequestMiddleware
}

func NewServer(config config.Config, srvc *service.Service) *Server {
	return &Server{
		application: newApplication(srvc),
		auth:        newAuth(srvc),
		system:      newSystem(srvc),
		user:        newUser(srvc),
		price:       newPrice(srvc),

		authMiddleware:    middleware.NewAuthMiddleware(srvc),
		timeoutMiddleware: middleware.NewTimeoutRequestMiddleware(time.Duration(config.HTTP.RequestTimeout) * time.Second),
	}
}

func customOutgoingHeaderMatcher(h string) (string, bool) {
	return h, true
}

func (s *Server) Register(router *mux.Router) {
	commonMiddlewares := []mux.MiddlewareFunc{
		middleware.LoggerMiddleware,
		middleware.CorsMiddleware,
	}
	privateMiddlewares := []mux.MiddlewareFunc{
		s.authMiddleware.RequestMiddleware,
	}
	managementMiddlewares := []mux.MiddlewareFunc{
		s.authMiddleware.RequestMiddleware,
		middleware.ManagementMiddleware,
	}

	ctx := context.Background()
	// Create grpc mux
	grpcMux := runtime.NewServeMux(
		// Fix "Grpc-Metadata-" prefix for outgoing headers
		runtime.WithOutgoingHeaderMatcher(customOutgoingHeaderMatcher),
	)

	// Init grpcMux with all APIs
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
	err = s.system.RegisterHTTP(ctx, grpcMux)
	if err != nil {
		logger.Error.Fatal("error register system", err.Error())
	}
	err = s.price.RegisterHTTP(ctx, grpcMux)
	if err != nil {
		logger.Error.Fatal("error register price", err.Error())
	}

	// management
	managementRoute := router.PathPrefix("/management").Subrouter()
	managementRoute.Use(managementMiddlewares...)
	managementRoute.PathPrefix("").Handler(grpcMux)

	// private
	privateRoute := router.PathPrefix("/private").Subrouter()
	privateRoute.Use(privateMiddlewares...)
	privateRoute.PathPrefix("").Handler(grpcMux)

	router.Use(commonMiddlewares...)
	router.PathPrefix("").Handler(grpcMux)
}
