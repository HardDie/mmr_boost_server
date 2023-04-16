package server

import (
	"time"

	"github.com/gorilla/mux"

	"github.com/HardDie/mmr_boost_server/internal/config"
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

func (s *Server) Register(router *mux.Router) {
	privateMiddlewares := []mux.MiddlewareFunc{
		middleware.LoggerMiddleware,
		s.timeoutMiddleware.RequestMiddleware,
		s.authMiddleware.RequestMiddleware,
	}

	applicationsRouter := router.PathPrefix("/applications").Subrouter()
	s.application.RegisterPrivateRouter(applicationsRouter, privateMiddlewares...)

	authRouter := router.PathPrefix("/auth").Subrouter()
	s.auth.RegisterPublicRouter(authRouter, middleware.LoggerMiddleware)
	s.auth.RegisterPrivateRouter(authRouter, privateMiddlewares...)

	systemRouter := router.PathPrefix("/system").Subrouter()
	s.system.RegisterPublicRouter(systemRouter, middleware.LoggerMiddleware, middleware.CorsMiddleware, s.timeoutMiddleware.RequestMiddleware)

	userRouter := router.PathPrefix("/user").Subrouter()
	s.user.RegisterPrivateRouter(userRouter, privateMiddlewares...)
}
