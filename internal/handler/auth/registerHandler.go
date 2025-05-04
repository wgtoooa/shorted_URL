package auth

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"url_shortness/internal/repository/database"
	"url_shortness/internal/repository/database/query"
)

func RegisterHandlerGet(ctx *gin.Context) {
	ctx.HTML(200, "register.html", nil)
}

func RegisterHandlerPost(ctx *gin.Context) {
	login := strings.TrimSpace(ctx.PostForm("login"))
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(ctx.PostForm("password")), bcrypt.DefaultCost) // hash user password
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed hash password"})
		return
	}

	checkLogin, err := query.AccountExists(database.DB, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if checkLogin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "login already exists"})
		return
	}

	err = query.CreateAccount(database.DB, login, string(hashPassword))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed create account"})
		return
	}
	log.Println("successfully")
	ctx.Redirect(http.StatusFound, "/login")
}
