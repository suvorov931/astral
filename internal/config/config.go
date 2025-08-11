package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"astral/internal/api"
	"astral/internal/auth"
	"astral/internal/logger"
)

type Config struct {
	HttpServer api.HttpServer
	//Redis       redisClient.Config
	//Postgres    postgresClient.Config
	Auth   auth.Config
	Logger logger.Config
}

func New(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
