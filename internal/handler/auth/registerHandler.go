package auth

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"url_shortness/internal/repository/database/query"
)

func (a *Auth) RegisterHandlerGet(ctx *gin.Context) {
	ctx.HTML(200, "register.html", nil) // register form
}

func (a *Auth) RegisterHandlerPost(ctx *gin.Context) {
	login := ValidLogin(ctx.PostForm("login"))
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(ctx.PostForm("password")), bcrypt.DefaultCost) // hash user password
	if err != nil {
		a.logger.Error("failed hash password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed hash password"})
		return
	}

	checkLogin, err := query.AccountExists(a.db, login) // check login exist
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		a.logger.Error("failed check login")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if checkLogin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "login already exists"})
		return
	}

	err = query.CreateAccount(a.db, login, string(hashPassword)) // just create new account and add in database
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed create account"})
		return
	}
	a.logger.Info("successfully,account created")
	ctx.Redirect(http.StatusFound, "/login") // moving to the login page
}
