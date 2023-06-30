package application

import (
	"context"
	"fmt"
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
	"github.com/HardDie/mmr_boost_server/internal/repository/encrypt"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/repository/smtp"
	"github.com/HardDie/mmr_boost_server/internal/server"
	"github.com/HardDie/mmr_boost_server/internal/service"
)

const (
	ServerTimeout   = 30
	ShutdownTimeout = 30
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
		return nil, fmt.Errorf("error init db: %w", err)
	}
	app.DB = newDB

	// Init migrations
	err = migration.NewMigrate(app.DB).Up()
	if err != nil {
		return nil, fmt.Errorf("error run migrations: %w", err)
	}

	// Prepare router
	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()

	// Init repositories
	repositoryAccessToken := postgres.NewAccessToken(app.DB)
	repositoryApplication := postgres.NewApplication(app.DB)
	repositoryEmailValidation := postgres.NewEmailValidation(app.DB)
	repositoryResetPassword := postgres.NewResetPassword(app.DB)
	repositoryHistory := postgres.NewHistory(app.DB)
	repositoryPassword := postgres.NewPassword(app.DB)
	repositoryUser := postgres.NewUser(app.DB)
	repositoryStatusHistory := postgres.NewStatusHistory(app.DB)
	postgresRepository := postgres.NewPostgres(
		app.DB,
		repositoryAccessToken,
		repositoryApplication,
		repositoryEmailValidation,
		repositoryResetPassword,
		repositoryHistory,
		repositoryPassword,
		repositoryUser,
		repositoryStatusHistory,
	)
	smtpRepository := smtp.NewSMTP(app.Cfg)
	encryptRepository, err := encrypt.NewEncrypt(app.Cfg.Encrypt)
	if err != nil {
		return nil, fmt.Errorf("error init encrypt repo: %w", err)
	}

	// Init services
	serviceApplication := service.NewApplication(postgresRepository, encryptRepository)
	serviceAuth := service.NewAuth(app.Cfg, postgresRepository, smtpRepository)
	serviceSystem := service.NewSystem()
	serviceUser := service.NewUser(postgresRepository)
	servicePrice := service.NewPrice(postgresRepository)
	serviceStatusHistory := service.NewStatusHistory(postgresRepository)
	srvc := service.NewService(
		serviceApplication,
		serviceAuth,
		serviceSystem,
		serviceUser,
		servicePrice,
		serviceStatusHistory,
	)

	// Init severs
	srv := server.NewServer(app.Cfg, srvc)
	srv.Register(v1Router)

	return app, nil
}

func (app *Application) Run() error {
	srv := &http.Server{
		Addr:         app.Cfg.HTTP.Port,
		ReadTimeout:  ServerTimeout * time.Second,
		WriteTimeout: ServerTimeout * time.Second,
		Handler:      app.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	defer app.Stop()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Println("server error:", err.Error())
		}
	}()

	<-done
	logger.Info.Println("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout*time.Second)
	defer func() {
		cancel()
	}()

	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error.Println("error shutdown server:", err.Error())
		return err
	}
	return nil
}

func (app *Application) Stop() {
	app.DB.DB.Close()
	app.DB = nil
	logger.Info.Println("Done")
}
