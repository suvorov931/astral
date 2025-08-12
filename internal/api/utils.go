package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type mainResponse struct {
	ErrorResponse *ErrorResponse `json:"error,omitempty"`
	Response      *Response      `json:"response,omitempty"`
}

type ErrorResponse struct {
	Code int    `json:"code,omitempty"`
	Text string `json:"text,omitempty"`
}

func WriteError(w http.ResponseWriter, logger *zap.Logger, code int, text string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := mainResponse{
		ErrorResponse: &ErrorResponse{
			Code: code,
			Text: text,
		},
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("WriteError: failed to encode response", zap.Error(err))
	}
}

type Response struct {
	Login string `json:"login,omitempty"`
	Token string `json:"token,omitempty"`
}

func WriteResponseWithLogin(w http.ResponseWriter, logger *zap.Logger, login string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := mainResponse{
		Response: &Response{
			Login: login,
		},
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("WriteResponseWithLogin: failed to encode response", zap.Error(err))
	}
}

func WriteResponseWithToken(w http.ResponseWriter, logger *zap.Logger, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := mainResponse{
		Response: &Response{
			Token: token,
		},
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("WriteResponseWithToken: failed to encode response", zap.Error(err))
	}
}
