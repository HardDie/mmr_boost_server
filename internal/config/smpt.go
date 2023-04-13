package config

type SMTP struct {
	Host     string
	Port     int
	Email    string
	Password string
	Nickname string
}

func smtpConfig() SMTP {
	return SMTP{
		Host:     getEnv("SMTP_HOST"),
		Port:     getEnvAsInt("SMTP_PORT"),
		Email:    getEnv("SMTP_EMAIL"),
		Password: getEnv("SMTP_PASSWORD"),
		Nickname: getEnv("SMTP_NICKNAME"),
	}
}
