package service

import (
	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/repository/postgres"
	"github.com/HardDie/mmr_boost_server/internal/repository/smtp"
)

type Service struct {
	auth
	system
	user
}

func NewService(
	config config.Config,
	postgresRepository *postgres.Postgres,
	smtpRepository *smtp.SMTP,
) *Service {
	return &Service{
		auth:   newAuth(config, postgresRepository, smtpRepository),
		system: newSystem(),
		user:   newUser(postgresRepository),
	}
}
