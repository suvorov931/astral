package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"astral/internal/api"
	"astral/internal/auth"
)

const bearerPrefix = "Bearer "

func RequireAdminToken(as auth.AuthService, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header, err := getAuthorizationHeader(r)
			if err != nil {
				api.WriteError(w, logger, http.StatusUnauthorized, "No authorization header found")
				logger.Error("RequireAdminToken:", zap.Error(err))
				return
			}

			token, err := extractToken(header)
			if err != nil {
				api.WriteError(w, logger, http.StatusUnauthorized, "Invalid authorization header format")
				logger.Error("RequireAdminToken:", zap.Error(err))
				return
			}

			if isAdminToken := as.IsAdminToken(token); !isAdminToken {
				api.WriteError(w, logger, http.StatusUnauthorized, "Invalid admin token")
				logger.Error("RequireAdminToken: invalid admin token")
				return
			}

			logger.Info("RequireAdminToken: admin token is correctly")

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func getAuthorizationHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", fmt.Errorf("getAuthorizationHeader: no authorization header found")
	}

	return header, nil
}

func extractToken(header string) (string, error) {
	if !strings.HasPrefix(header, bearerPrefix) {
		return "", fmt.Errorf("extractToken: invalid authorization header format")
	}

	tokenString := strings.TrimPrefix(header, bearerPrefix)

	return tokenString, nil
}
