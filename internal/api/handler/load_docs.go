package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"astral/internal/api"
	"astral/internal/auth"
	redisClient "astral/internal/cache/redisCLient"
	"astral/internal/documents"
	"astral/internal/storage/postgresClient"
)

const maxLoadSize = 50 << 20

func LoadDocs(pc postgresClient.PostgresClient, rc redisClient.RedisClient, as auth.AuthService, logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r.Body = http.MaxBytesReader(w, r.Body, maxLoadSize)

		err := r.ParseMultipartForm(maxLoadSize)
		if err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, "invalid form data")
			logger.Warn("LoadDocs: invalid form data", zap.Error(err))
			return
		}

		metaStr := r.FormValue("meta")
		if metaStr == "" {
			api.WriteError(w, logger, http.StatusBadRequest, "meta required")
			logger.Warn("LoadDocs: meta is missing")
			return
		}

		var meta api.Meta

		err = json.Unmarshal([]byte(metaStr), &meta)
		if err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, "invalid meta json")
			logger.Warn("LoadDocs: invalid meta json", zap.Error(err))
			return
		}

		if meta.Name == "" && !meta.File {
			api.WriteError(w, logger, http.StatusBadRequest, "name required")
			logger.Warn("LoadDocs: name is missing")
			return
		}

		hashToken := as.GenerateSha(meta.Token)

		login, err := rc.GetLoginByToken(ctx, hashToken)
		if err != nil {
			api.WriteError(w, logger, http.StatusUnauthorized, "invalid token")
			logger.Warn("LoadDocs: invalid token", zap.Error(err))
			return
		}

		id := uuid.NewString()

		document := documents.Document{
			Id:        id,
			Login:     login,
			Name:      meta.Name,
			Mime:      meta.Mime,
			File:      meta.File,
			Public:    meta.Public,
			CreatedAt: time.Now(),
		}

		if jsonStr := r.FormValue("json"); jsonStr != "" {
			document.JSON = []byte(jsonStr)
		}

		if !meta.Public && len(meta.Grant) == 0 {
			document.Grant = []string{login}
		}

		if meta.File {
			file, header, err := r.FormFile("file")
			if err != nil {
				api.WriteError(w, logger, http.StatusBadRequest, "file is required")
				logger.Warn("LoadDocs: file is required", zap.Error(err))
				return
			}
			defer file.Close()

			if document.Name == "" {
				document.Name = filepath.Base(header.Filename)
			}

			content, err := io.ReadAll(file)
			if err != nil {
				api.WriteError(w, logger, http.StatusInternalServerError, "failed to read file")
				logger.Error("LoadDocs: failed read file", zap.Error(err))
				return
			}

			document.Content = content
		}

		err = pc.SaveDocument(ctx, &document)
		if err != nil {
			api.WriteError(w, logger, http.StatusInternalServerError, "failed to save document")
			logger.Error("LoadDocs: failed save document", zap.Error(err))
			return
		}

		//if err := rc.InvalidateUserDocsList(ctx, document.Login); err != nil {
		//	logger.Warn("LoadDocs: failed to invalidate docs list cache", zap.Error(err))
		//}
		//if err := rc.InvalidateDoc(ctx, document.Id); err != nil {
		//	logger.Warn("LoadDocs: failed to invalidate doc cache", zap.Error(err))
		//}
		//
		var jsonData interface{}
		if len(document.JSON) > 0 {
			if err := json.Unmarshal(document.JSON, &jsonData); err != nil {
				jsonData = string(document.JSON)
			}
		} else {
			jsonData = nil
		}

		api.WriteResponseWithData(w, logger, jsonData, document.Name)
		logger.Info("LoadDocs: successfully loaded document", zap.String("id", id))
	}
}
