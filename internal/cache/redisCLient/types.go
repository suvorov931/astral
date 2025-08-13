package redisClient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type Config struct {
	Host     string        `env:"REDIS_HOST" env-required:"true"`
	Port     int           `env:"REDIS_PORT" env-required:"true"`
	TokenDB  int           `env:"REDIS_TOKEN_DB" env-required:"true"`
	CacheDB  int           `env:"REDIS_CACHE_DB" env-required:"true"`
	Timeout  time.Duration `env:"REDIS_TIMEOUT" env-required:"true"`
	TokenTTL time.Duration `env:"REDIS_TOKEN_TTL" env-required:"true"`
	CacheTTL time.Duration `env:"REDIS_CACHE_TTL" env-required:"true"`
	Password string        `env:"REDIS_PASSWORD" env-required:"true"`
}

type RedisService struct {
	tokenDB  *redis.Client
	cacheDB  *redis.Client
	logger   *zap.Logger
	timeout  time.Duration
	tokenTTL time.Duration
	cacheTTL time.Duration
}

type RedisClient interface {
	SaveToken(ctx context.Context, key string, token string) error
	Close()
}

type MockRedisClient struct {
	mock.Mock
}
