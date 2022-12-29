package config

type AppConfig struct {
	AppPort string `env:"APP_PORT" default:"3000"`
}
