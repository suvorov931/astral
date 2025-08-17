package handler

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"astral/internal/api"
	"astral/internal/auth"
	"astral/internal/storage/postgres_client"
	"astral/internal/storage/redis_client"
)

// Auth godoc
// @Summary      Authenticate user and return token
// @Description  Validate user credentials and return generated token. Token is stored server-side (redis) as hashed value.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      api.User  true  "User credentials (login + password)"
// @Success      200   {object}  api.mainResponse  "Returns generated token in response"
// @Failure      400   {object}  api.mainResponse  "Invalid request body"
// @Failure      401   {object}  api.mainResponse  "User not found or invalid password"
// @Failure      500   {object}  api.mainResponse  "Server error (DB/Redis/Token generation)"
// @Router       /api/auth [post]
func Auth(pc postgresClient.PostgresClient, rc redisClient.RedisClient, as auth.AuthService, logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var user api.User

		err := decodeBody(w, r, &user)
		if err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, "invalid request body")
			logger.Warn("Auth: invalid request body", zap.Error(err))
			return
		}

		storedHash, err := pc.GetPasswordHash(ctx, user.Login)
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

		tokenHash := as.GenerateSha(token)

		err = rc.SaveToken(ctx, user.Login, tokenHash)
		if err != nil {
			api.WriteError(w, logger, http.StatusInternalServerError, "cannot save token")
			logger.Error("Auth: cannot save token", zap.Error(err))
			return
		}

		api.WriteResponseWithToken(w, logger, token)
		logger.Info("Auth: successfully validate user, generate and save token")
	}
}
