package auth

import "go.uber.org/zap"

type Config struct {
	AdminToken string `env:"ADMIN_TOKEN"`
}

type Auth struct {
	config *Config
	logger *zap.Logger
}

type AuthService interface {
	IsAdminToken(token string) bool
}
