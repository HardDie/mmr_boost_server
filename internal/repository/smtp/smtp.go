package smtp

import (
	"fmt"

	"github.com/go-mail/mail/v2"

	"github.com/HardDie/mmr_boost_server/internal/config"
)

const (
	Subject = "E-mail verification"
)

type SMTP struct {
	email    string
	nickname string
	dialer   *mail.Dialer
}

func NewSMTP(cfg config.SMTP) *SMTP {
	return &SMTP{
		email:    cfg.Email,
		nickname: cfg.Nickname,
		dialer:   mail.NewDialer(cfg.Host, cfg.Port, cfg.Email, cfg.Password),
	}
}

func (r *SMTP) SendEmailVerification(email, code string) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email)
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", r.nickname, r.email))
	msg.SetHeader("Subject", Subject)
	msg.SetBody("text/plain", fmt.Sprintf("Your code: %s", code))

	err := r.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("smtp: error sending message: %w", err)
	}
	return nil
}
