package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"

	"astral/internal/logger"
)

type Config struct {
	//HttpServer  api.HttpServer
	//SMTP        SMTPClient.Config
	//Redis       redisClient.Config
	//Postgres    postgresClient.Config
	Logger logger.Config
}

func New(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
