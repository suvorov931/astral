package auth

import (
	"crypto/subtle"

	"go.uber.org/zap"
)

func New(config *Config, logger *zap.Logger) *Auth {
	return &Auth{
		config: config,
		logger: logger,
	}
}

func (a *Auth) IsAdminToken(token string) bool {
	if subtle.ConstantTimeCompare([]byte(token), []byte(a.config.AdminToken)) == 1 {
		return true
	} else {
		return false
	}
}
