package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/HardDie/mmr_boost_server/internal/logger"
)

var (
	env string
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type Config struct {
	App             App
	Postgres        Postgres
	HTTP            HTTP
	Password        Password
	SMTP            SMTP
	Session         Session
	EmailValidation EmailValidation
	Encrypt         Encrypt
}

func Get() Config {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Fatalf("failed to load env vars: %s", err)
		}
	}

	cfg := Config{
		App:             appConfig(),
		Postgres:        postgresConfig(),
		HTTP:            httpConfig(),
		Password:        passwordConfig(),
		SMTP:            smtpConfig(),
		Session:         sessionConfig(),
		EmailValidation: emailValidationConfig(),
		Encrypt:         encryptConfig(),
	}
	env = cfg.App.Env
	return cfg
}

func GetEnv() string {
	return env
}
