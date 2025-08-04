package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
	"url_shortness/internal/infra/product"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

//var jwtSecret = []byte("YOUR_SECRET_KEY")

func (a *Auth) LoginHandlerPost(ctx *gin.Context) {
	var inputLoginDate struct { // it information user send from login page
		Login    string `form:"login"`
		Password string `form:"password"`
	}

	if err := ctx.ShouldBind(&inputLoginDate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка получения данных"}) //get information
		return
	}

	login := ValidLogin(inputLoginDate.Login)

	user, err := product.GetUser(a.db, a.rdb, login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if !CheckPassword(user.Password, inputLoginDate.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильный пароль"})
		return
	}

	token, err := GenerateToken(inputLoginDate.Login) // generate token in order to user verification
	if err != nil {
		a.logger.Error("Ошибка генерации токена", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
		return
	}
	a.logger.Info("token", zap.String("token", token), zap.String("login", user.Login))
	ctx.SetCookie("token", token, 3600, "/", "", false, false) // set Cookie

	// И ОБЯЗАТЕЛЬНО вернуть токен в теле ответа
	ctx.JSON(200, gin.H{
		"token":    token,
		"redirect": "/protected/url", // Куда перейти после входа
	})
}

func (a *Auth) LoginHandlerGet(ctx *gin.Context) {
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
