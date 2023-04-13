package config

type Postgres struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func postgresConfig() Postgres {
	return Postgres{
		Host:     getEnv("POSTGRES_HOST"),
		Port:     getEnvAsInt("POSTGRES_PORT"),
		User:     getEnv("POSTGRES_USER"),
		Password: getEnv("POSTGRES_PASSWORD"),
		Database: getEnv("POSTGRES_DB"),
	}
}
