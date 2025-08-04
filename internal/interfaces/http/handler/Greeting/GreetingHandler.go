package Greeting

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
)

type Greeting struct {
	rdb    *redis.Client
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewGreeting(pool *pgxpool.Pool, logger *zap.Logger, redis *redis.Client) *Greeting {
	return &Greeting{
		rdb:    redis,
		db:     pool,
		logger: logger,
	}
}

func (h *Greeting) GreetingHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "greeting.html", nil)
}
