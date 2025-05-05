package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	handler "url_shortness/internal/handler"
	"url_shortness/internal/handler/auth"
	"url_shortness/internal/handler/shortURL"
	"url_shortness/internal/repository/database"
	"url_shortness/internal/repository/database/query"
)

const DSN = "host=dpg-d0cgip6uk2gs73981jeg-a user=wgtoooa password=j217PeAIMy9KrwFoFIV9ftYwIBuk0OXs dbname=shorted_url port=5432 sslmode=disable" //"postgres://postgres:12345@localhost:5432/postgres"

func main() {
	router := gin.Default()

	database.BDinit(DSN) //initialize database
	defer database.DB.Close()

	err := query.CreatedTableAccount(database.DB) // create table accounts if exists
	if err != nil {
		log.Printf("failsed create table account %e", err)
		return
	}
	err = query.CreateTableURL(database.DB) // create table URL if exists
	if err != nil {
		log.Printf("failsed create table URL %e", err)
		return
	}

	files, err := filepath.Glob("templates/**/*") // load all files from templates
	if err != nil {
		log.Println("failed in filePath")
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
