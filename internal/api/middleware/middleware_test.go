package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"astral/internal/auth"
)

func TestRequireAdminToken(t *testing.T) {
	as := auth.New(&auth.Config{
		AdminToken: "someAdminToken",
	}, zap.NewNop())

	tests := []struct {
		name       string
		token      string
		statusCode int
		response   string
	}{
		{
			name:       "valid token",
			token:      bearerPrefix + "someAdminToken",
			statusCode: http.StatusOK,
			response:   http.StatusText(http.StatusOK),
		},
		{
			name:       "invalid token",
			token:      bearerPrefix + "wrongToken",
			statusCode: http.StatusUnauthorized,
			response:   "{\"error\":{\"code\":401,\"text\":\"Invalid admin token\"}}\n",
		},
		{
			name:       "invalid header format",
			token:      "wrongToken",
			statusCode: http.StatusUnauthorized,
			response:   "{\"error\":{\"code\":401,\"text\":\"Invalid authorization header format\"}}\n",
		},
		{
			name:       "empty header",
			token:      "",
			statusCode: http.StatusUnauthorized,
			response:   "{\"error\":{\"code\":401,\"text\":\"No authorization header found\"}}\n",
		},
		{
			name:       "empty token",
			token:      bearerPrefix,
			statusCode: http.StatusUnauthorized,
			response:   "{\"error\":{\"code\":401,\"text\":\"Invalid admin token\"}}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(http.StatusText(http.StatusOK)))
			})

			handlerToTest := RequireAdminToken(as, zap.NewNop())(nextHandler)

			r := httptest.NewRequest("POST", "/api/register", nil)
			r.Header.Set("Authorization", tt.token)
			w := httptest.NewRecorder()

			handlerToTest.ServeHTTP(w, r)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.response, w.Body.String())
		})
	}
}
