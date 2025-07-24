package shortURL

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type URL struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewShortURL(pool *pgxpool.Pool, logger *zap.Logger) *URL {
	return &URL{
		db:     pool,
		logger: logger,
	}
}
