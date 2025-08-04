package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"url_shortness/internal/infra/database/query"
)

func (a *Auth) RegisterHandlerGet(ctx *gin.Context) {
	ctx.HTML(200, "register.html", nil) // register form
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Auth) RegisterHandlerPost(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	login := ValidLogin(req.Login)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost) // hash user password
	if err != nil {
		a.logger.Error("failed hash password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	checkLogin, err := query.AccountExists(a.db, login) // check login exist
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		a.logger.Error("Ошибка проверки логина", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}
	if checkLogin {
		a.logger.Warn("Попытка регистрации существующего логина",
			zap.String("login", login),
		)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "такой login уже существует"})
		return
	}

	err = query.CreateAccount(a.db, login, string(hashPassword)) // just create new account and add in database
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"redirect": "/login", // Опционально
		"message":  "Registration successful",
	})
}
