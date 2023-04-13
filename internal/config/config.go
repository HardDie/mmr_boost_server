package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/HardDie/mmr_boost_server/internal/logger"
)

type Config struct {
	Postgres Postgres
	Http     Http
	Password Password
	SMTP     SMTP
	Session  Session
}

func Get() Config {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Fatalf("failed to load env vars: %s", err)
		}
	}

	return Config{
		Postgres: postgresConfig(),
		Http:     httpConfig(),
		Password: passwordConfig(),
		SMTP:     smtpConfig(),
		Session:  sessionConfig(),
	}
}
