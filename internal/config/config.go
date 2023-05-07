package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/HardDie/mmr_boost_server/internal/logger"
)

type Config struct {
	Postgres        Postgres
	HTTP            HTTP
	Password        Password
	SMTP            SMTP
	Session         Session
	EmailValidation EmailValidation
}

func Get() Config {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Fatalf("failed to load env vars: %s", err)
		}
	}

	return Config{
		Postgres:        postgresConfig(),
		HTTP:            httpConfig(),
		Password:        passwordConfig(),
		SMTP:            smtpConfig(),
		Session:         sessionConfig(),
		EmailValidation: emailValidationConfig(),
	}
}
