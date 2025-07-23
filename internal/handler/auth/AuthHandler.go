package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Auth struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewAuth(pool *pgxpool.Pool, logger *zap.Logger) *Auth {
	return &Auth{
		db:     pool,
		logger: logger,
	}
}
