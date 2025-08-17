package handler

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"astral/internal/api"
	"astral/internal/auth"
	"astral/internal/storage/postgres_client"
)

const sizeLimit = 1 << 20

// Register godoc
// @Summary      Create a new user
// @Description  Create a new user. This endpoint is protected by an admin token middleware.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      api.User  true  "User credentials (login + password)"
// @Success      200   {object}  api.mainResponse  "Returns created login"
// @Failure      400   {object}  api.mainResponse  "Validation error or duplicate login"
// @Failure      500   {object}  api.mainResponse  "Server error"
// @Security     BearerAuth
// @Router       /api/register [post]
func Register(ps postgresClient.PostgresClient, as auth.AuthService, logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var user api.User

		err := decodeBody(w, r, &user)
		if err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, "invalid request body")
			logger.Warn("Register: invalid request body", zap.Error(err))
			return
		}

		err = as.ValidateLogin(user.Login)
		if err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, processValidateError(err))
			logger.Warn("Register:", zap.Error(err))
			return
		}

		err = as.ValidatePassword(user.Pswd)
		if err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, processValidateError(err))
			logger.Warn("Register:", zap.Error(err))
			return
		}

		passwordHash, err := as.HashPassword(user.Pswd)
		if err != nil {
			api.WriteError(w, logger, http.StatusInternalServerError, "cannot hash password")
			logger.Warn("Register:", zap.Error(err))
			return
		}

		err = ps.SaveUser(ctx, user.Login, string(passwordHash))
		if err != nil {
			switch {
			case errors.Is(err, postgresClient.ErrDuplicateLogin):
				api.WriteError(w, logger, http.StatusBadRequest, "duplicate login")
				logger.Error("Register:", zap.Error(err))
				return

			default:
				api.WriteError(w, logger, http.StatusInternalServerError, "cannot save user")
				logger.Error("Register:", zap.Error(err))
				return
			}
		}

		api.WriteResponseWithLogin(w, logger, user.Login)

		logger.Info("Register: successfully create and save user")
	}
}

func processValidateError(err error) string {
	switch {
	case errors.Is(err, auth.ErrShortLogin):
		return "login must be at least 8 characters"
	case errors.Is(err, auth.ErrInvalidLogin):
		return "invalid login"

	case errors.Is(err, auth.ErrShortPassword):
		return "password must be at least 8 characters"
	case errors.Is(err, auth.ErrMissingUpper):
		return "password must contain at least one uppercase letter"
	case errors.Is(err, auth.ErrMissingLower):
		return "password must contain at least one lowercase letter"
	case errors.Is(err, auth.ErrMissingDigit):
		return "password must contain at least one digit"
	case errors.Is(err, auth.ErrMissingSpecial):
		return "password must contain at least one special symbol"

	default:
		return "invalid input"
	}
}
