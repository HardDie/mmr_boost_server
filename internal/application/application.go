package application

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/migration"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/repository/smtp"
	"github.com/HardDie/mmr_boost_server/internal/server"
	"github.com/HardDie/mmr_boost_server/internal/service"
)

type Application struct {
	Cfg    config.Config
	DB     *db.DB
	Router *mux.Router
}

func Get() (*Application, error) {
	app := &Application{
		Cfg:    config.Get(),
		Router: mux.NewRouter(),
	}

	// Init DB
	newDB, err := db.Get(app.Cfg.Postgres)
	if err != nil {
		return nil, err
	}
	app.DB = newDB

	// Init migrations
	err = migration.NewMigrate(app.DB).Up()
	if err != nil {
		return nil, err
	}

	// Prepare router
	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()

	// Init repositories
	postgresRepository := postgres.NewPostgres(app.DB)
	smtpRepository := smtp.NewSMTP(app.Cfg.SMTP)

	// Init services
	srvc := service.NewService(app.Cfg, postgresRepository, smtpRepository)

	// Init severs
	srv := server.NewServer(app.Cfg, srvc)
	srv.Register(v1Router)

	return app, nil
}

func (app *Application) Run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		app.Stop()
		os.Exit(0)
	}()

	defer app.Stop()
	return http.ListenAndServe(app.Cfg.Http.Port, app.Router)
}

func (app *Application) Stop() {
	app.DB.DB.Close()
	app.DB = nil
	log.Println("Done")
}
