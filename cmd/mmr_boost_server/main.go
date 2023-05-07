// Package main MMR Boost Sever.
//
// Entry point for the application.
//
// Terms Of Service:
//
//	Schemes: https
//	Host: event-tracker.online
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- token:
//
//	SecurityDefinitions:
//	token:
//	     type: apiKey
//	     name: Authorization
//	     in: header
//
// swagger:meta
package main

import (
	"github.com/HardDie/mmr_boost_server/internal/application"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

func main() {
	app, err := application.Get()
	if err != nil {
		logger.Error.Fatal(err)
	}
	logger.Info.Println("Server listen on", app.Cfg.HTTP.Port)
	err = app.Run()
	if err != nil {
		logger.Error.Fatal(err)
	}
}
