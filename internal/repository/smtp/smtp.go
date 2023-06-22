package smtp

import (
	"bytes"
	"fmt"

	"github.com/go-mail/mail/v2"

	"github.com/HardDie/mmr_boost_server/internal/config"
	"github.com/HardDie/mmr_boost_server/internal/logger"
)

const (
	Subject = "E-mail verification"
)

type SMTP struct {
	email    string
	nickname string
	baseURL  string
	dialer   *mail.Dialer
}

func NewSMTP(cfg config.Config) *SMTP {
	return &SMTP{
		email:    cfg.SMTP.Email,
		nickname: cfg.SMTP.Nickname,
		baseURL:  cfg.EmailValidation.BaseURL,
		dialer:   mail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Email, cfg.SMTP.Password),
	}
}

func (r *SMTP) SendEmailVerification(email, code string) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email)
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", r.nickname, r.email))
	msg.SetHeader("Subject", Subject)
	msg.SetBody("text/plain",
		fmt.Sprintf("Your code: %s, or link %s/api/v1/auth/validate_email?code=%s", code, r.baseURL, code),
	)
	msg.AddAlternative("text/html",
		fmt.Sprintf(
			"<p>Your code: <b>%s</b>, or link <a href=%s/api/v1/auth/validate_email?code=%s>validate email</a></p>",
			code, r.baseURL, code,
		),
	)

	if config.GetEnv() == config.ENV_DEV {
		b := bytes.NewBuffer(nil)
		_, err := msg.WriteTo(b)
		if err != nil {
			return fmt.Errorf("smtp: error write message to buffer: %w", err)
		}
		logger.Debug.Println("Email:", b.String())
		return nil
	}

	err := r.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("smtp: error sending message: %w", err)
	}
	return nil
}
