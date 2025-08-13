package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
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
	if len(login) < loginLength {
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
	if len(password) < passwordLength {
		a.logger.Warn("ValidatePassword:", zap.Error(ErrShortPassword))
		return ErrShortPassword
	}

	var upper, lower, digit, special bool

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			upper = true
		case char >= 'a' && char <= 'z':
			lower = true
		case char >= '0' && char <= '9':
			digit = true
		default:
			special = true
		}
	}

	switch {
	case !upper:
		a.logger.Warn("ValidatePassword:", zap.Error(ErrMissingUpper))
		return ErrMissingUpper

	case !lower:
		a.logger.Warn("ValidatePassword:", zap.Error(ErrMissingLower))
		return ErrMissingLower

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

func (a *Auth) GenerateToken() (string, error) {
	buf := make([]byte, a.config.LengthToken)

	_, err := rand.Read(buf)
	if err != nil {
		a.logger.Error("GenerateToken: failed to generate token", zap.Error(err))
		return "", fmt.Errorf("GenerateToken: failed to generate token: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(buf)

	return token, nil
}

func (a *Auth) GenerateSha(token string) string {
	shaHash := sha256.Sum256([]byte(token))

	hashString := hex.EncodeToString(shaHash[:])

	return hashString
}
