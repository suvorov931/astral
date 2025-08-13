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
	HTTP_PORT=4848

	ADMIN_TOKEN=someAdminToken
	LENGTH_TOKEN=111

	REDIS_HOST=localhost
	REDIS_PORT=6754321
	REDIS_TOKEN_DB=0
	REDIS_CACHE_DB=1
	REDIS_TIMEOUT=33s
	REDIS_PASSWORD=redisPassword

	POSTGRES_HOST=localhost
	POSTGRES_PORT=6754321
	POSTGRES_USER=root
	POSTGRES_PASSWORD=postgresPassword
	POSTGRES_DATABASE=postgres
	POSTGRES_TIMEOUT=33s
	POSTGRES_MAX_CONNECTIONS=1000
	POSTGRES_MIN_CONNECTIONS=500

	LOGGER=dev
	`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := New(tempFile)
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.HttpServer.Host)
	assert.Equal(t, 4848, cfg.HttpServer.Port)

	assert.Equal(t, "someAdminToken", cfg.Auth.AdminToken)
	assert.Equal(t, 111, cfg.Auth.LengthToken)

	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6754321, cfg.Redis.Port)
	assert.Equal(t, 0, cfg.Redis.TokenDB)
	assert.Equal(t, 1, cfg.Redis.CacheDB)
	assert.Equal(t, 33*time.Second, cfg.Redis.Timeout)
	assert.Equal(t, "redisPassword", cfg.Redis.Password)

	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, "6754321", cfg.Postgres.Port)
	assert.Equal(t, "root", cfg.Postgres.User)
	assert.Equal(t, "postgresPassword", cfg.Postgres.Password)
	assert.Equal(t, "postgres", cfg.Postgres.Database)
	assert.Equal(t, 33*time.Second, cfg.Postgres.Timeout)
	assert.Equal(t, 1000, cfg.Postgres.MaxConns)
	assert.Equal(t, 500, cfg.Postgres.MinConns)

	assert.Equal(t, "dev", cfg.Logger.Env)

	_, err = New("wrongPath")
	assert.Contains(t, err.Error(), "failed to read config")
}
