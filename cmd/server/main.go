package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	handler "url_shortness/internal/handler"
	"url_shortness/internal/handler/auth"
	"url_shortness/internal/handler/shortURL"
	"url_shortness/internal/repository/database"
	"url_shortness/pkg/logger"
)

var DSlN = database.DataBaseConfig{
	User:     os.Getenv("DB_USER"),
	Password: os.Getenv("DB_PASSWORD"),
	Host:     os.Getenv("DB_HOST"),
	Port:     os.Getenv("DB_PORT"),
	DBName:   os.Getenv("DB_NAME"),
}

var DSN = database.DataBaseConfig{
	User:     "admin",
	Password: "secret",
	Host:     "localhost",
	Port:     "5432",
	DBName:   "mydb",
}

func main() {

	router := gin.Default()

	log.Println(DSN)
	database.BDinit(DSN) //initialize database
	defer database.DB.Close()

	files, err := filepath.Glob("templates/**/*") // load all files from templates
	if err != nil {
		logger.Log.Error("failed Load template files", zap.Error(err))
		return
	}

	router.LoadHTMLFiles(files...) // load all files in server

	router.GET("/", handler.HelloHandler)

	router.GET("/register", auth.RegisterHandlerGet)
	router.POST("/register", auth.RegisterHandlerPost)

	router.GET("/login", auth.LoginHandlerGet)
	router.POST("/login", auth.LoginHandlerPost)

	Protected := router.Group("/")
	Protected.Use(auth.AuthMiddleware())

	Protected.GET("/url", shortURL.ShowURLhandler)
	Protected.GET("/url/data", shortURL.GetUserURLsJSON)
	Protected.POST("/url", shortURL.CreateShortURLhandler)
	Protected.GET("/:url", shortURL.FollowURLHandler)

	log.Println("server starting....")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run("0.0.0.0:" + port)
}

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Error("Error loading .env file")
		return
	}
}
