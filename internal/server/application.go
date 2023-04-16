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

type application struct {
	service *service.Service
}

func newApplication(service *service.Service) application {
	return application{
		service: service,
	}
}

func (s *application) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	applicationRouter := router.PathPrefix("").Subrouter()
	applicationRouter.HandleFunc("", s.ApplicationCreate).Methods(http.MethodPost)
	applicationRouter.Use(middleware...)
}

// swagger:parameters ApplicationCreateRequest
type ApplicationCreateRequest struct {
	// In: body
	Body struct {
		dto.ApplicationCreateRequest
	}
}

// swagger:response ApplicationCreateResponse
type ApplicationCreateResponse struct {
	// In: body
	Body struct {
		Data *entity.ApplicationPublic `json:"data"`
	}
}

// swagger:route POST /api/v1/applications Application ApplicationCreateRequest
//
// # Create application for boosting
//
//	Responses:
//	  201: ApplicationCreateResponse
func (s *application) ApplicationCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &dto.ApplicationCreateRequest{}
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

	resp, err := s.service.ApplicationCreate(ctx, req)
	if err != nil {
		errs.HttpError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = utils.Response(w, resp)
	if err != nil {
		logger.Error.Println("error write to socket:", err.Error())
	}
}
