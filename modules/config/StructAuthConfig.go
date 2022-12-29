package config

type AuthConfig struct {
	JwtSecret string `env:"JWT_SECRET" required:"true"`
}
