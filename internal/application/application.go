package application

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/db"
	"github.com/HardDie/mmr_boost_server/internal/logger"
	"github.com/HardDie/mmr_boost_server/internal/migration"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/repository/smtp"
	"github.com/HardDie/mmr_boost_server/internal/server"
	"github.com/HardDie/mmr_boost_server/internal/service"
)

const (
	ServerTimeout = 30
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
	repositoryAccessToken := postgres.NewAccessToken(app.DB)
	repositoryApplication := postgres.NewApplication(app.DB)
	repositoryEmailValidation := postgres.NewEmailValidation(app.DB)
	repositoryHistory := postgres.NewHistory(app.DB)
	repositoryPassword := postgres.NewPassword(app.DB)
	repositoryUser := postgres.NewUser(app.DB)
	postgresRepository := postgres.NewPostgres(
		app.DB,
		repositoryAccessToken,
		repositoryApplication,
		repositoryEmailValidation,
		repositoryHistory,
		repositoryPassword,
		repositoryUser,
	)
	smtpRepository := smtp.NewSMTP(app.Cfg)

	// Init services
	serviceApplication := service.NewApplication(postgresRepository)
	serviceAuth := service.NewAuth(app.Cfg, postgresRepository, smtpRepository)
	serviceSystem := service.NewSystem()
	serviceUser := service.NewUser(postgresRepository)
	servicePrice := service.NewPrice(postgresRepository)
	srvc := service.NewService(
		serviceApplication,
		serviceAuth,
		serviceSystem,
		serviceUser,
		servicePrice,
	)

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

	srv := &http.Server{
		Addr:         app.Cfg.HTTP.Port,
		ReadTimeout:  ServerTimeout * time.Second,
		WriteTimeout: ServerTimeout * time.Second,
		Handler:      app.Router,
	}
	return srv.ListenAndServe()
}

func (app *Application) Stop() {
	app.DB.DB.Close()
	app.DB = nil
	logger.Info.Println("Done")
}
