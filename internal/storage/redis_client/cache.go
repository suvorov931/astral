package redisClient

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"astral/internal/documents"
)

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

	rs.logger.Info("CacheDocument: completed cache", zap.String("doc", document.Id))
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
