package config

import (
	"github.com/joho/godotenv"
	"github.com/yousefggg/common-lib/pkg/config"
	"github.com/yousefggg/common-lib/pkg/logger"
)

type Config struct {
	*config.Config
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		logger.Info(".env file not found, using system environment variables")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		panic(err) 
	}

	return &Config{
		Config: cfg,
	}
}