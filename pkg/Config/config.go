package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"time"
	"url_shortness/internal/infra/redis"
	"url_shortness/pkg/logger"
)

type DataBaseConfig struct {
	User            string        `envconfig:"POSTGRES_USER" required:"true" `
	DBName          string        `envconfig:"POSTGRES_NAME" required:"true" `
	Password        string        `envconfig:"POSTGRES_PASSWORD" required:"true" `
	Host            string        `envconfig:"POSTGRES_HOST" required:"true" `
	Port            string        `envconfig:"POSTGRES_PORT" required:"true"`
	SSLMode         string        `envconfig:"DB_SSL_MODE" default:"disable"`
	MaxConns        int32         `envconfig:"DB_MAX_CONNS" default:"10"`
	MinConns        int32         `envconfig:"DB_MIN_CONNS" default:"2"`
	MaxConnLifetime time.Duration `envconfig:"DB_MAX_CONN_LIFETIME" default:"1h"`
	MaxConnIdleTime time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME" default:"30m"`
	ConnectTimeout  time.Duration `envconfig:"DB_CONNECT_TIMEOUT" default:"5s"`
}
type Config struct {
	Production     bool `envconfig:"PRODUCTION" default:"false"`
	DataBaseConfig DataBaseConfig
	Server         struct {
		Port string `envconfig:"SERVER_PORT" default:"8080"`
		Host string `envconfig:"SERVER_HOST" default:"localhost"`
	}
	ConfigRedis redis.ConfigRedis
}

func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Get().Info("No .env file found, using environment variables")
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Get().Fatal("Failed to process env vars", zap.Error(err))
		return nil, err
	}
	fmt.Println(cfg)
	return &cfg, nil
}
