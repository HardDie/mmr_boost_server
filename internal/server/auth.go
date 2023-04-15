package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/mmr_boost_server/internal/dto"
	"github.com/HardDie/mmr_boost_server/internal/entity"
	"github.com/HardDie/mmr_boost_server/internal/errs"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/service"
	"github.com/HardDie/mmr_boost_server/internal/utils"
)

type auth struct {
	service *service.Service
}

func newAuth(service *service.Service) auth {
	return auth{
		service: service,
	}
}
func (s *auth) RegisterPublicRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.HandleFunc("/register", s.Register).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", s.Login).Methods(http.MethodPost)
	authRouter.Use(middleware...)
}
func (s *auth) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.HandleFunc("/user", s.User).Methods(http.MethodGet)
	authRouter.HandleFunc("/logout", s.Logout).Methods(http.MethodPost)
	authRouter.Use(middleware...)
}

/*
 * Public
 */

// swagger:parameters AuthRegisterRequest
type AuthRegisterRequest struct {
	// In: body
	Body struct {
		dto.AuthRegisterRequest
	}
}

// swagger:response AuthRegisterResponse
type AuthRegisterResponse struct {
	// In: body
	Body struct {
		Data *entity.AccessToken `json:"data"`
	}
}

// swagger:route POST /api/v1/auth/register Auth AuthRegisterRequest
//
// # Registration form
//
//	Responses:
//	  200: AuthRegisterResponse
func (s *auth) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &dto.AuthRegisterRequest{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = getValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.service.AuthRegister(ctx, req)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	accessToken, err := s.service.AuthGenerateCookie(ctx, user.ID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	utils.SetSessionCookie(accessToken.TokenHash, w)

	w.WriteHeader(http.StatusCreated)
	err = utils.Response(w, accessToken)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}

// swagger:parameters AuthLoginRequest
type AuthLoginRequest struct {
	// In: body
	Body struct {
		dto.AuthLoginRequest
	}
}

// swagger:response AuthLoginResponse
type AuthLoginResponse struct {
	// In: body
	Body struct {
		Data *entity.AccessToken `json:"data"`
	}
}

// swagger:route POST /api/v1/auth/login Auth AuthLoginRequest
//
// # Login form
//
//	Responses:
//	  200: AuthLoginResponse
func (s *auth) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &dto.AuthLoginRequest{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = getValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.service.AuthLogin(ctx, req)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	accessToken, err := s.service.AuthGenerateCookie(ctx, user.ID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	utils.SetSessionCookie(accessToken.TokenHash, w)

	err = utils.Response(w, accessToken)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}

/*
 * Private
 */

// swagger:parameters AuthUserRequest
type AuthUserRequest struct {
}

// swagger:response AuthUserResponse
type AuthUserResponse struct {
	// In: body
	Body struct {
		Data *entity.User `json:"data"`
	}
}

// swagger:route GET /api/v1/auth/user Auth AuthUserRequest
//
// # Getting information about the current user
//
//	Responses:
//	  200: AuthUserResponse
func (s *auth) User(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	user, err := s.service.AuthGetUserInfo(ctx, userID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	err = utils.Response(w, user)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}

// swagger:parameters AuthLogoutRequest
type AuthLogoutRequest struct {
}

// swagger:response AuthLogoutResponse
type AuthLogoutResponse struct {
}

// swagger:route POST /api/v1/auth/logout Auth AuthLogoutRequest
//
// # Close the current session
//
//	Responses:
//	  200: AuthLogoutResponse
func (s *auth) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := utils.GetAccessTokenFromContext(ctx)

	err := s.service.AuthLogout(ctx, session.ID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	utils.DeleteSessionCookie(w)
}
