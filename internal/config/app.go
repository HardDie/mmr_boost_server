package config

type App struct {
	Env string
}

func appConfig() App {
	return App{
		Env: getEnv("APP_ENV"),
	}
}
