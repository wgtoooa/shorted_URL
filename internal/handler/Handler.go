package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewHandler(pool *pgxpool.Pool, logger *zap.Logger) *Handler {
	return &Handler{
		db:     pool,
		logger: logger,
	}
}
