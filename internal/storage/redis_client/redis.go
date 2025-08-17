package redisClient

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func New(ctx context.Context, config *Config, logger *zap.Logger) (*RedisService, error) {
	tokenDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.TokenDB,
	})
	err := tokenDB.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	cacheDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.CacheDB,
	})
	err = cacheDB.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &RedisService{
		tokenDB:  tokenDB,
		cacheDB:  cacheDB,
		logger:   logger,
		timeout:  config.Timeout,
		tokenTTL: config.TokenTTL,
		cacheTTL: config.CacheTTL,
	}, nil
}

func (rs *RedisService) Close() {
	err := rs.cacheDB.Close()
	if err != nil {
		rs.logger.Warn("Close: failed to close cacheDB", zap.Error(err))
	}

	err = rs.tokenDB.Close()
	if err != nil {
		rs.logger.Warn("Close: failed to close tokenDB", zap.Error(err))
	}
}
