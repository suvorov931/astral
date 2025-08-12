package auth

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

const (
	loginLength    = 8
	passwordLength = 8

	loginPattern = `^[A-Za-z0-9]+$`
)

var (
	ErrShortLogin   = errors.New("short login")
	ErrInvalidLogin = errors.New("invalid login")

	ErrShortPassword  = errors.New("short password")
	ErrMissingUpper   = errors.New("miss uppercase letter")
	ErrMissingLower   = errors.New("miss lowercase letter")
	ErrMissingDigit   = errors.New("miss digit")
	ErrMissingSpecial = errors.New("miss special symbol")
)

type Config struct {
	AdminToken  string `env:"ADMIN_TOKEN"`
	LengthToken int    `env:"LENGTH_TOKEN"`
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
	GenerateToken() (string, error)
}

type MockAuthService struct {
	mock.Mock
}
