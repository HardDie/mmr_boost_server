package config

type Session struct {
	AccessToken  int
	RefreshToken int
}

func sessionConfig() Session {
	return Session{
		AccessToken:  getEnvAsInt("SESSION_ACCESS_TOKEN"),
		RefreshToken: getEnvAsInt("SESSION_REFRESH_TOKEN"),
	}
}
