package Services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"url_shortness/internal/handler/Greeting"
	"url_shortness/internal/handler/auth"
	"url_shortness/internal/handler/shortURL"
	logger2 "url_shortness/internal/logger"
)

type AppServices interface {
	Handler() *Greeting.Greeting
	URL() *shortURL.URL
	Auth() *auth.Auth
}
type appServices struct {
	logger  *zap.Logger
	pool    *pgxpool.Pool
	handler *Greeting.Greeting
	url     *shortURL.URL
	auth    *auth.Auth
}

func NewAppServices(pool *pgxpool.Pool) *appServices {
	logger := logger2.Get()
	return &appServices{
		pool:    pool,
		logger:  logger,
		handler: Greeting.NewGreeting(pool, logger),
		url:     shortURL.NewShortURL(pool, logger),
		auth:    auth.NewAuth(pool, logger),
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
