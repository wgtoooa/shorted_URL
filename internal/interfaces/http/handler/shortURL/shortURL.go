package shortURL

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type URL struct {
	rdb    *redis.Client
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewShortURL(pool *pgxpool.Pool, logger *zap.Logger, redis *redis.Client) *URL {
	return &URL{
		rdb:    redis,
		db:     pool,
		logger: logger,
	}
}
