package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"url_shortness/pkg/logger"
)

type ConfigRedis struct {
	Addr     string `envconfig:"REDIS_ADDRESS" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:"*****"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
	PoolSize int    `envconfig:"REDIS_POOL_SIZE" default:"10"`
	//	DefaultTTL time.Duration `envconfig:"REDIS_DEFAULT_TTL" default:"10m"`
}

func InitRedis(cfg ConfigRedis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Get().Panic(err.Error())
	}
	logger.Get().Info("Redis service initialized")

	return client
}
