package handler

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"url_shortness/internal/repository/Table"
	"url_shortness/internal/repository/database"
	"url_shortness/internal/repository/database/query"
)

func LoginHandlerPost(ctx *gin.Context) {
	login := ctx.PostForm("login")
	password := strings.TrimSpace(ctx.PostForm("password"))

	user, err := query.GetAccount(database.DB, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "such account does no exist"})
			return
			return
		}

		log.Printf("database error when fetching login %s %v ", login, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if user == (Table.Account{}) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "such account does no exist"})
		return
	}

	if !CheckPassword(user.Password, password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}
	ctx.Redirect(http.StatusFound, "/url")
}

func LoginHandlerGet(ctx *gin.Context) {
	ctx.HTML(200, "login.html", nil)
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
