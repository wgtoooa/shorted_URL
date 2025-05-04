package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
	handler "url_shortness/internal/handler"
	"url_shortness/internal/repository/database"
	"url_shortness/internal/repository/database/query"
)

var DSN = "postgres://postgres:12345@localhost:5432/postgres"

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

	router.GET("/register", handler.RegisterHandlerGet)
	router.POST("/register", handler.RegisterHandlerPost)

	router.GET("/login", handler.LoginHandlerGet)
	router.POST("/login", handler.LoginHandlerPost)

	router.GET("/url", handler.ShowURLhandler)

	log.Println("server starting....")
	router.Run(":8080") //starting server
}
