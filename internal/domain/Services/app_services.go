package Services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"url_shortness/internal/interfaces/http/handler/Greeting"
	"url_shortness/internal/interfaces/http/handler/auth"
	"url_shortness/internal/interfaces/http/handler/shortURL"
	logger2 "url_shortness/pkg/logger"
)

type AppServices interface {
	Handler() *Greeting.Greeting
	URL() *shortURL.URL
	Auth() *auth.Auth
}
type appServices struct {
	redis   *redis.Client
	logger  *zap.Logger
	pool    *pgxpool.Pool
	handler *Greeting.Greeting
	url     *shortURL.URL
	auth    *auth.Auth
}

func NewAppServices(pool *pgxpool.Pool, rdb *redis.Client) *appServices {
	logger := logger2.Get()
	return &appServices{
		redis:   rdb,
		pool:    pool,
		logger:  logger,
		handler: Greeting.NewGreeting(pool, logger, rdb),
		url:     shortURL.NewShortURL(pool, logger, rdb),
		auth:    auth.NewAuth(pool, logger, rdb),
	}
}

func (s *appServices) Handler() *Greeting.Greeting {
	return s.handler
}

func (s *appServices) URL() *shortURL.URL {
	return s.url
}

func (s *appServices) Auth() *auth.Auth {
	return s.auth
}
