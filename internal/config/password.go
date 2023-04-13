package config

type Password struct {
	MaxAttempts int `json:"maxAttempts"`
	BlockTime   int `json:"blockTime"`
}

func passwordConfig() Password {
	return Password{
		MaxAttempts: getEnvAsInt("PASSWORD_MAX_ATTEMPTS"),
		BlockTime:   getEnvAsInt("PASSWORD_BLOCK_TIME"),
	}
}
