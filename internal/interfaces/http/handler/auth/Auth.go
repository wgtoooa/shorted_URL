package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strings"
)

type Auth struct {
	rdb    *redis.Client
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAuth(pool *pgxpool.Pool, logger *zap.Logger, redis *redis.Client) *Auth {
	return &Auth{
		rdb:    redis,
		db:     pool,
		logger: logger,
	}
}

func ValidLogin(login string) string {
	Login := strings.TrimSpace(login)
	Login = strings.ToLower(Login)
	return Login
}
