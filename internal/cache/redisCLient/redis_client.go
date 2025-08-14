package redisClient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"astral/internal/documents"
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

func (rs *RedisService) SaveToken(ctx context.Context, login string, token string) error {
	ctx, cancel := context.WithTimeout(ctx, rs.timeout)
	defer cancel()

	err := rs.tokenDB.Set(ctx, token, login, rs.tokenTTL).Err()
	if err != nil {
		rs.logger.Error("SaveToken: failed to save token", zap.Error(err))
		return fmt.Errorf("SaveToken: failed to save token: %w", err)
	}

	rs.logger.Info("SaveToken: successfully saved token")
	return nil
}

func (rs *RedisService) GetLoginByToken(ctx context.Context, token string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, rs.timeout)
	defer cancel()

	login, err := rs.tokenDB.Get(ctx, token).Result()
	if err != nil {
		rs.logger.Error("GetLoginByToken: failed to get token", zap.Error(err))
		return "", fmt.Errorf("GetLoginByToken: failed to get token: %w", err)
	}

	return login, nil
}

func (rs *RedisService) CacheDocument(ctx context.Context, document *documents.Document) error {
	ctx, cancel := context.WithTimeout(ctx, rs.timeout)
	defer cancel()

	docForCache := *document
	docForCache.Content = nil

	docBytes, err := json.Marshal(docForCache)
	if err != nil {
		rs.logger.Warn("CacheDocument: failed to marshal document for cache", zap.Error(err))
		return fmt.Errorf("CacheDocument: failed to marshal document for cache: %w", err)
	}

	err = rs.cacheDB.Set(ctx, "doc:"+docForCache.Id, docBytes, rs.cacheTTL).Err()
	if err != nil {
		rs.logger.Warn("CacheDocument: failed to cache document", zap.Error(err))
	}

	pattern := "docs:" + document.Login + ":*"

	var cursor uint64

	for {
		keys, next, err := rs.cacheDB.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			rs.logger.Warn("CacheDocument: scan failed for docs list invalidation", zap.Error(err), zap.String("pattern", pattern))
			break
		}

		for i := 0; i < len(keys); i += batchSize {
			end := i + batchSize

			if end > len(keys) {
				end = len(keys)
			}

			chunk := keys[i:end]
			err = rs.cacheDB.Del(ctx, chunk...).Err()
			if err != nil {
				rs.logger.Warn("CacheDocument: failed to del keys chunk", zap.Error(err), zap.Int("chunk_size", len(chunk)))

			} else {
				rs.logger.Debug("CacheDocument: deleted keys chunk", zap.Int("chunk_size", len(chunk)))
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	rs.logger.Info("CacheDocument: completed cache+invalidation", zap.String("doc", document.Id))
	return nil
}

func (rs *RedisService) InvalidateDocs(ctx context.Context, login string) error {
	ctx, cancel := context.WithTimeout(ctx, rs.timeout)
	defer cancel()

	pattern := "docs:" + login + ":*"
	var cursor uint64

	for {
		keys, next, err := rs.cacheDB.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			rs.logger.Warn("InvalidateDocs: scan failed", zap.Error(err), zap.String("pattern", pattern))
			return err
		}

		for i := 0; i < len(keys); i += batchSize {
			end := i + batchSize
			if end > len(keys) {
				end = len(keys)
			}

			chunk := keys[i:end]
			err = rs.cacheDB.Del(ctx, chunk...).Err()
			if err != nil {
				rs.logger.Warn("InvalidateDocs: del failed for chunk", zap.Error(err), zap.Int("chunk_size", len(chunk)))

			} else {
				rs.logger.Debug("InvalidateDocs: deleted keys chunk", zap.Int("chunk_size", len(chunk)))
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	rs.logger.Info("InvalidateDocs: successfully invalidated docs for login", zap.String("login", login))
	return nil
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
