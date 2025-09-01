package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"url_shortness/internal/domain/Services"
	"url_shortness/internal/infra/database"
	"url_shortness/internal/infra/redis"
	"url_shortness/internal/interfaces/http/handler/NotFound"
	"url_shortness/pkg/Config"
	"url_shortness/pkg/logger"

	"path/filepath"
)

func main() {

	logger.Init(false) //--Сделать еще config
	defer logger.Get().Sync()

	router := gin.Default()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Get().Fatal("Error loading config", zap.Error(err))
	}

	var DSN = cfg.DataBaseConfig

	dbPool, err := database.InitDB(DSN)
	if err != nil {
		logger.Get().Fatal("Database initialization failed", zap.Error(err))
	}
	defer dbPool.Close()

	rdb := redis.InitRedis(cfg.ConfigRedis)
	defer rdb.Close()

	services := Services.NewAppServices(dbPool, rdb)

	files, err := filepath.Glob("templates/**/*") // load all files from templates
	if err != nil {
		logger.Get().Error("failed Load template files", zap.Error(err))
		return
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	router.LoadHTMLFiles(files...) // load all files in server
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	logger.Get().Info("templates and static loaded ")

	router.GET("/", services.Handler().GreetingHandler)

	router.GET("/register", services.Auth().RegisterHandlerGet)
	router.POST("/register", services.Auth().RegisterHandlerPost)

	router.GET("/login", services.Auth().LoginHandlerGet)
	router.POST("/login", services.Auth().LoginHandlerPost)

	Protected := router.Group("/protected")
	Protected.Use(services.Auth().AuthMiddleware())

	Protected.POST("/url", services.URL().CreateShortURLHandler)
	Protected.GET("/url", services.URL().ShowURLHandler)
	Protected.PATCH("/url", services.URL().PatchURLHandler)
	Protected.GET("/url/data", services.URL().GetUserURLsJSON)
	Protected.GET("/l/:url", services.URL().FollowURLHandler)
	Protected.DELETE("/url/:short_url", services.URL().DeleteURL)

	router.NoRoute(NotFound.NotFoundHandler)

	cfgServer := cfg.Server

	err = router.Run(cfgServer.Host + ":" + cfgServer.Port)
	if err != nil {
		logger.Get().Fatal("Error running server", zap.Error(err))
	}
	logger.Get().Info("server starting....", zap.String("port", cfgServer.Port))
}
