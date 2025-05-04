package auth

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"url_shortness/internal/repository/database"
	"url_shortness/internal/repository/database/query"
)


func LoginHandlerPost(ctx *gin.Context) {
	var creds struct {
		Login    string `form:"login"`
		Password string `form:"password"`
	}

	if err := ctx.ShouldBind(&creds); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed get data about creds"})
		return
	}

	user, err := query.GetAccount(database.DB, creds.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "such account does not exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	if !CheckPassword(user.Password, creds.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}

	token, err := GenerateToken(creds.Login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
	ctx.Redirect(http.StatusFound, "/url")
}

func LoginHandlerGet(ctx *gin.Context) {
	ctx.HTML(200, "login.html", nil)
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(login string) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
