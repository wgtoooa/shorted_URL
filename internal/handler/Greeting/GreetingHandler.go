package Greeting

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
)

type Greeting struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewGreeting(pool *pgxpool.Pool, logger *zap.Logger) *Greeting {
	return &Greeting{
		db:     pool,
		logger: logger,
	}
}

func (h *Greeting) GreetingHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "greeting.html", nil)
}
