package auth

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

const (
	loginPattern = `^[A-Za-z0-9]+$`
)

var (
	ErrShortLogin   = errors.New("short login")
	ErrInvalidLogin = errors.New("invalid login")

	ErrShortPassword    = errors.New("short password")
	ErrMissingUpper     = errors.New("miss uppercase letter")
	ErrMissingLower     = errors.New("miss lowercase letter")
	ErrMissingDigit     = errors.New("miss digit")
	ErrMissingSpecial   = errors.New("miss special symbol")
	ErrMissingMixedCase = errors.New("miss mixed case")
)

type Config struct {
	AdminToken string `env:"ADMIN_TOKEN"`
}

type Auth struct {
	config *Config
	logger *zap.Logger
}

type AuthService interface {
	IsAdminToken(token string) bool
	ValidateLogin(login string) error
	ValidatePassword(password string) error
	HashPassword(password string) ([]byte, error)
}

type MockAuthService struct {
	mock.Mock
}
