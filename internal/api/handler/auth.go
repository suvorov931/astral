package handler

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"astral/internal/api"
	"astral/internal/auth"
	"astral/internal/storage/postgresClient"
)

func Auth(ps postgresClient.PostgresClient, as auth.AuthService, logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var user api.User

		err := decodeBody(w, r, &user)
		if err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, "invalid request body")
			logger.Warn("Auth: invalid request body", zap.Error(err))
			return
		}

		storedHash, err := ps.GetPasswordHash(ctx, user.Login)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				api.WriteError(w, logger, http.StatusUnauthorized, "user not found")
				logger.Warn("Auth: user not found")
				return
			}

			api.WriteError(w, logger, http.StatusInternalServerError, "cannot get stored hash")
			logger.Error("Auth: cannot get stored hash", zap.Error(err))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Pswd))
		if err != nil {
			api.WriteError(w, logger, http.StatusUnauthorized, "invalid password")
			logger.Warn("Auth: invalid user")
			return
		}

		token, err := as.GenerateToken()
		if err != nil {
			api.WriteError(w, logger, http.StatusInternalServerError, "cannot generate token")
			logger.Warn("Auth: cannot generate token", zap.Error(err))
			return
		}

		api.WriteResponseWithToken(w, logger, token)
		logger.Info("Auth: successfully validate user and generate token")
	}
}

// TODO: somehow store hash
