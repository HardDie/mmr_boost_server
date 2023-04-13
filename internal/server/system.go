package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/service"
)

type system struct {
	service *service.Service
}

func newSystem(service *service.Service) system {
	return system{
		service: service,
	}
}

func (s *system) RegisterPublicRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	systemRouter := router.PathPrefix("").Subrouter()
	systemRouter.HandleFunc("/swagger", s.Swagger).Methods(http.MethodGet)
	systemRouter.Use(middleware...)
}

/*
 * Public
 */

// swagger:parameters SwaggerRequest
type SwaggerRequest struct {
}

// swagger:response SwaggerResponse
type SwaggerResponse struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/v1/system/swagger System SwaggerRequest
//
// # Get the yaml-file of the swagger description
//
//	Responses:
//	  200: SwaggerResponse
func (s *system) Swagger(w http.ResponseWriter, r *http.Request) {
	data, err := s.service.SystemGetSwagger()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(data)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
