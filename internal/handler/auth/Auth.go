package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
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

func ValidLogin(login string) string {
	Login := strings.TrimSpace(login)
	Login = strings.ToLower(Login)
	return Login
}
