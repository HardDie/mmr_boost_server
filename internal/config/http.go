package config

type Http struct {
	Port           string `json:"port"`
	RequestTimeout int    `json:"requestTimeout"`
}

func httpConfig() Http {
	return Http{
		Port:           getEnv("HTTP_PORT"),
		RequestTimeout: getEnvAsInt("HTTP_REQUEST_TIMEOUT"),
	}
}
