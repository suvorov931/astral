package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"astral/internal/api"
	"astral/internal/auth"
	"astral/internal/logger"
	"astral/internal/storage/postgres_client"
	"astral/internal/storage/redis_client"
)

type Config struct {
	HttpServer api.HttpServer        `env-required:"true"`
	Auth       auth.Config           `env-required:"true"`
	Redis      redisClient.Config    `env-required:"true"`
	Postgres   postgresClient.Config `env-required:"true"`
	Logger     logger.Config         `env-required:"true"`
}

func New(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
