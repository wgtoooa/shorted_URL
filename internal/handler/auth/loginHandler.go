package auth

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
	"url_shortness/internal/repository/database/query"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

//var jwtSecret = []byte("YOUR_SECRET_KEY")

func (h *Auth) LoginHandlerPost(ctx *gin.Context) {
	var inputLoginDate struct { // it information user send from login page
		Login    string `form:"login"`
		Password string `form:"password"`
	}

	if err := ctx.ShouldBind(&inputLoginDate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed get data about inputLoginDate"}) //get information
		return
	}

	user, err := query.GetAccount(h.db, inputLoginDate.Login) // get account it user
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "such account does not exist"})
			return
		}
		h.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	if !CheckPassword(user.Password, inputLoginDate.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}

	token, err := GenerateToken(inputLoginDate.Login) // generate token in order to user verification
	if err != nil {
		h.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	h.logger.Info("token", zap.String("token", token), zap.String("login", user.Login))
	ctx.SetCookie("token", token, 3600, "/", "localhost", true, true) // set Cookie
	ctx.Redirect(http.StatusFound, "/url")
}

func (h *Auth) LoginHandlerGet(ctx *gin.Context) {
	ctx.HTML(200, "login.html", nil)
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(login string) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
