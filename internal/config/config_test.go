package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := tempDir + "/config.env"

	content := `
	HTTP_HOST=localhost
	HTTP_PORT=8080
	HTTP_MONITORING_PORT=2112
	HTTP_TIMEOUT_EXTRA=3s

	ADMIN_TOKEN=someAdminToken
	LENGTH_TOKEN=11

	POSTGRES_HOST=localhost
	POSTGRES_PORT=5432
	POSTGRES_USER=root
	POSTGRES_PASSWORD=postgresPassword
	POSTGRES_DATABASE=postgres
	POSTGRES_TIMEOUT=3s
	POSTGRES_MAX_CONNECTIONS=10
	POSTGRES_MIN_CONNECTIONS=5

	LOGGER=dev
	`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := New(tempFile)
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.HttpServer.Host)
	assert.Equal(t, 8080, cfg.HttpServer.Port)

	assert.Equal(t, "someAdminToken", cfg.Auth.AdminToken)
	assert.Equal(t, 11, cfg.Auth.LengthToken)

	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, "5432", cfg.Postgres.Port)
	assert.Equal(t, "root", cfg.Postgres.User)
	assert.Equal(t, "postgresPassword", cfg.Postgres.Password)
	assert.Equal(t, "postgres", cfg.Postgres.Database)
	assert.Equal(t, 3*time.Second, cfg.Postgres.Timeout)
	assert.Equal(t, 10, cfg.Postgres.MaxConns)
	assert.Equal(t, 5, cfg.Postgres.MinConns)

	assert.Equal(t, "dev", cfg.Logger.Env)

	_, err = New("wrongPath")
	assert.Contains(t, err.Error(), "failed to read config")
}
