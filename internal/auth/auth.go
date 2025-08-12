package auth

import (
	"crypto/subtle"
	"fmt"
	"regexp"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func (a *Auth) ValidateLogin(login string) error {
	if len(login) < 8 {
		a.logger.Warn("ValidateLogin:", zap.Error(ErrShortLogin))
		return ErrShortLogin
	}

	match, err := regexp.MatchString(loginPattern, login)
	if err != nil {
		a.logger.Warn("ValidateLogin: cannot validate login", zap.Error(err))
		return fmt.Errorf("ValidateLogin: cannot validate login: %w", err)
	}
	if !match {
		a.logger.Warn("ValidateLogin:", zap.Error(ErrInvalidLogin))
		return ErrInvalidLogin
	}

	return nil
}

func (a *Auth) ValidatePassword(password string) error {
	if len(password) < 8 {
		a.logger.Warn("ValidatePassword:", zap.Error(ErrShortPassword))
		return ErrShortPassword
	}

	var upper, lower int
	var digit, special bool

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			upper++
		case char >= 'a' && char <= 'z':
			lower++
		case char >= '0' && char <= '9':
			digit = true
		default:
			special = true
		}
	}

	switch {
	case upper < 0:
		a.logger.Warn("ValidatePassword:", zap.Error(ErrMissingUpper))
		return ErrMissingUpper

	case lower < 0:
		a.logger.Warn("ValidatePassword:", zap.Error(ErrMissingLower))
		return ErrMissingLower

	case upper+lower < 2:
		a.logger.Warn("ValidatePassword:", zap.Error(ErrMissingMixedCase))
		return ErrMissingMixedCase

	case !digit:
		a.logger.Warn("ValidatePassword:", zap.Error(ErrMissingDigit))
		return ErrMissingDigit

	case !special:
		a.logger.Warn("ValidatePassword:", zap.Error(ErrMissingSpecial))
		return ErrMissingSpecial

	}

	return nil
}

func (a *Auth) HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("HashPassword: failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("HashPassword: failed to hash password: %w", err)
	}

	return hash, nil
}
