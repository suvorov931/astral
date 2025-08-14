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
	"astral/internal/storage/postgresClient"
)

const maxUploadSize = 50 << 20

//type Meta struct {
//	Name   string   `json:"name"`
//	File   bool     `json:"file"`
//	Public bool     `json:"public"`
//	Token  string   `json:"token"`
//	Mime   string   `json:"mime"`
//	Grant  []string `json:"grant"`
//}
//
//type Document struct {
//	Id        string
//	Login     string
//	Name      string
//	Mime      string
//	File      bool
//	Public    bool
//	Grant     []string
//	Content   []byte
//	JSON      []byte
//	CreatedAt time.Time
//}

func LoadDocs(pc postgresClient.PostgresClient, rc redisClient.RedisClient, as auth.AuthService, logger *zap.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, "invalid form data")
			logger.Warn("LoadDocs: invalid form data", zap.Error(err))
			return
		}

		metaStr := r.FormValue("meta")
		if metaStr == "" {
			api.WriteError(w, logger, http.StatusBadRequest, "meta required")
			logger.Warn("LoadDocs: meta required")
			return
		}

		var meta api.Meta

		if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
			api.WriteError(w, logger, http.StatusBadRequest, "invalid meta json")
			logger.Warn("LoadDocs: invalid meta json", zap.Error(err))
			return
		}

		hashToken := as.GenerateSha(meta.Token)

		login, err := rc.GetLoginByToken(ctx, hashToken)
		if err != nil {
			api.WriteError(w, logger, http.StatusUnauthorized, "invalid token")
			logger.Warn("LoadDocs: invalid token", zap.Error(err))
			return
		}

		if meta.Name == "" && !meta.File {
			api.WriteError(w, logger, http.StatusBadRequest, "name required")
			logger.Warn("LoadDocs: name is missing", zap.Error(err))
			return
		}

		id := uuid.NewString()

		document := api.Document{
			Id:        id,
			Login:     login,
			Name:      meta.Name,
			Mime:      meta.Mime,
			File:      meta.File,
			Public:    meta.Public,
			Grant:     meta.Grant,
			CreatedAt: time.Now(),
		}

		if jsonStr := r.FormValue("json"); jsonStr != "" {
			document.JSON = []byte(jsonStr)
		}

		err = pc.SaveDocument(ctx, &document)
		if err != nil {
			api.WriteError(w, logger, http.StatusInternalServerError, "failed to save document")
			logger.Error("LoadDocs: failed save document", zap.Error(err))
			return
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

			//err = pc.SaveDocument(ctx, &document)
			//if err != nil {
			//	api.WriteError(w, logger, http.StatusInternalServerError, "failed to save document")
			//	logger.Warn("LoadDocs: failed save document", zap.Error(err))
			//	return
			//}
			//
			//if err := rc.InvalidateUserDocsList(ctx, document.Login); err != nil {
			//	logger.Warn("LoadDocs: failed to invalidate docs list cache", zap.Error(err))
			//}
			//if err := rc.InvalidateDoc(ctx, document.Id); err != nil {
			//	logger.Warn("LoadDocs: failed to invalidate doc cache", zap.Error(err))
			//}
			//
			//var jsonData interface{}
			//if len(document.JSON) > 0 {
			//	if err := json.Unmarshal(document.JSON, &jsonData); err != nil {
			//		jsonData = string(document.JSON)
			//	}
			//} else {
			//	jsonData = nil
			//}
			//
			//resp := make(map[string]interface{})
			//data := make(map[string]interface{})
			//if jsonData != nil {
			//	data["json"] = jsonData
			//}
			//if document.File {
			//	data["file"] = document.Name
			//}
			//resp["data"] = data
			//
			//w.Header().Set("Content-Type", "application/json")
			//w.WriteHeader(http.StatusOK)
			//if err := json.NewEncoder(w).Encode(resp); err != nil {
			//	logger.Error("LoadDocs: failed to encode response", zap.Error(err))
			//}
		}
	}
}
