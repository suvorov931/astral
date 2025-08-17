package redisClient

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

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
