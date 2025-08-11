package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type mainResponse struct {
	errorResponse `json:"error"`
}

type errorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func WriteError(w http.ResponseWriter, logger *zap.Logger, code int, text string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := mainResponse{
		errorResponse: errorResponse{
			Code: code,
			Text: text,
		},
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("WriteError: failed to encode response", zap.Error(err))
	}
}
