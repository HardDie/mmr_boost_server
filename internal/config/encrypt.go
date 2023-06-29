package config

type Encrypt struct {
	Key string `json:"key"`
}

func encryptConfig() Encrypt {
	return Encrypt{
		Key: getEnv("ENCRYPT_KEY"),
	}
}
