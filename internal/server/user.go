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

type user struct {
	service *service.Service
}

func newUser(service *service.Service) user {
	return user{
		service: service,
	}
}
func (s *user) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	userRouter := router.PathPrefix("").Subrouter()
	userRouter.HandleFunc("/password", s.Password).Methods(http.MethodPatch)
	userRouter.HandleFunc("/steam_id", s.UpdateSteamID).Methods(http.MethodPatch)
	userRouter.Use(middleware...)
}

/*
 * Private
 */

// swagger:parameters UserPasswordRequest
type UserPasswordRequest struct {
	// In: body
	Body struct {
		dto.UserUpdatePasswordRequest
	}
}

// swagger:response UserPasswordResponse
type UserPasswordResponse struct {
}

// swagger:route PATCH /api/v1/user/password User UserPasswordRequest
//
// # Updating the password for a user
//
//	Responses:
//	  200: UserPasswordResponse
func (s *user) Password(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.UserUpdatePasswordRequest{}
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

	err = s.service.UserPassword(ctx, req, userID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}
}

// swagger:parameters UserUpdateSteamIDRequest
type UserUpdateSteamIDRequest struct {
	// In: body
	Body struct {
		dto.UserUpdateSteamIDRequest
	}
}

// swagger:response UserUpdateSteamIDResponse
type UserUpdateSteamIDResponse struct {
	// In: body
	Body struct {
		Data *entity.User `json:"data"`
	}
}

// swagger:route PATCH /api/v1/user/steam_id User UserUpdateSteamIDRequest
//
// # Updating the steam id for a user
//
//	Responses:
//	  200: UserUpdateSteamIDResponse
func (s *user) UpdateSteamID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.UserUpdateSteamIDRequest{}
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

	u, err := s.service.UserUpdateSteamID(ctx, req, userID)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	err = utils.Response(w, u)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}
