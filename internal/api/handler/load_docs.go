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
	"astral/internal/documents"
	"astral/internal/storage/postgres_client"
	"astral/internal/storage/redis_client"
)

const maxLoadSize = 50 << 20

// LoadDocs godoc
// @Summary      Upload or create a document
// @Description  Upload a document (file or JSON). The request is multipart/form-data.
// @Tags         docs
// @Accept       multipart/form-data
// @Produce      json
// @Param        meta  formData  string  true   "JSON string with metadata. Example: {\"name\":\"file.txt\",\"file\":true,\"public\":false,\"token\":\"...\",\"mime\":\"text/plain\",\"grant\":[\"user1\"]}"
// @Param        file  formData  file    false  "File to upload (required if meta.file is true)"
// @Param        json  formData  string  false  "Optional JSON payload (when not uploading a binary file)"
// @Success      200   {object}  api.mainResponse  "Returns document JSON (if any) and file name"
// @Failure      400   {object}  api.mainResponse  "Invalid form data / missing meta / missing file"
// @Failure      401   {object}  api.mainResponse  "Invalid token"
// @Failure      500   {object}  api.mainResponse  "Server error (DB/Redis/IO)"
// @Router       /api/docs [post]
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

		document := &documents.Document{
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

		err = pc.SaveDocument(ctx, document)
		if err != nil {
			api.WriteError(w, logger, http.StatusInternalServerError, "failed to save document")
			logger.Error("LoadDocs: failed save document", zap.Error(err))
			return
		}

		err = rc.CacheDocument(ctx, document)
		if err != nil {
			logger.Warn("LoadDocs: failed to cache document", zap.Error(err))
		}

		err = rc.InvalidateDocs(ctx, document.Login)
		if err != nil {
			logger.Warn("LoadDocs: failed to invalidate doc cache", zap.Error(err))
		}

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
