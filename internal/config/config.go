package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/yousefggg/common-lib/pkg/logger" 
)

type Config struct {
	HTTPPort  string
	Env       string
	DBURL     string
	JWTSecret string
	JWTTTL    time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		logger.Info("info: .env file not found, using system environment variables")
	}

	ttlStr := os.Getenv("JWT_TTL")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to parse JWT_TTL '%s', falling back to default 24h", ttlStr), "error", err)
		ttl = time.Hour * 24
	}

	return &Config{
		HTTPPort:  os.Getenv("HTTP_PORT"),
		Env:       os.Getenv("ENV"),
		DBURL:     os.ExpandEnv(os.Getenv("DB_URL")), 
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTTTL:    ttl,
	}
}