package postgresClient

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type Config struct {
	Host     string        `env:"POSTGRES_HOST" env-required:"true"`
	Port     string        `env:"POSTGRES_PORT" env-required:"true"`
	User     string        `env:"POSTGRES_USER" env-required:"true"`
	Password string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database string        `env:"POSTGRES_DATABASE" env-required:"true"`
	Timeout  time.Duration `env:"POSTGRES_TIMEOUT" env-required:"true"`
	MaxConns int           `env:"POSTGRES_MAX_CONNECTIONS" env-required:"true"`
	MinConns int           `env:"POSTGRES_MIN_CONNECTIONS" env-required:"true"`
}

var (
	ErrDuplicateLogin = errors.New("duplicate login")
)

type PostgresService struct {
	pool    *pgxpool.Pool
	logger  *zap.Logger
	timeout time.Duration
}

type PostgresClient interface {
	SaveUser(ctx context.Context, login string, passwordHash string) error
	GetPasswordHash(ctx context.Context, login string) (string, error)
	Close()
}

type MockPostgresService struct {
	mock.Mock
}
