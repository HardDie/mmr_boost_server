package config

type EmailValidation struct {
	Expiration int    `json:"expiration"`
	BaseURL    string `json:"baseUrl"`
}

func emailValidationConfig() EmailValidation {
	return EmailValidation{
		Expiration: getEnvAsInt("EMAIL_VALIDATION_EXPIRATION"),
		BaseURL:    getEnv("EMAIL_VALIDATION_BASE_URL"),
	}
}
