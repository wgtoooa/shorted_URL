package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"url_shortness/internal/Services"
	"url_shortness/internal/logger"
	"url_shortness/internal/repository/database"
)

func main() {

	router := gin.Default()

	var DSN = database.DataBaseConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}
	logger.Init(false) //--Сделать еще config
	defer logger.Get().Sync()

	dbPool, err := database.InitDB(DSN)
	if err != nil {
		logger.Get().Fatal("Database initialization failed", zap.Error(err))
	}
	defer dbPool.Close()

	services := Services.NewAppServices(dbPool)

	files, err := filepath.Glob("templates/**/*") // load all files from templates
	if err != nil {
		logger.Get().Error("failed Load template files", zap.Error(err))
		return
	}

	router.LoadHTMLFiles(files...) // load all files in server

	logger.Get().Info("templates loaded")

	router.GET("/", services.Handler().GreetingHandler)

	router.GET("/register", services.Auth().RegisterHandlerGet)
	router.POST("/register", services.Auth().RegisterHandlerPost)

	router.GET("/login", services.Auth().LoginHandlerGet)
	router.POST("/login", services.Auth().LoginHandlerPost)

	Protected := router.Group("/")
	Protected.Use(services.Auth().AuthMiddleware())

	Protected.GET("/url", services.URL().ShowURLHandler)
	Protected.GET("/url/data", services.URL().GetUserURLsJSON)
	Protected.POST("/url", services.URL().CreateShortURLHandler)
	Protected.GET("/:url", services.URL().FollowURLHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run("0.0.0.0:" + port)
	logger.Get().Info("server starting....", zap.String("port", port))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}
}
