package shortURL

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/idna"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"url_shortness/internal/handler/auth"
	"url_shortness/internal/repository/database/Table"
	"url_shortness/internal/repository/database/query"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var SeedeRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func (u *URL) ShowURLHandler(ctx *gin.Context) {
	login := ctx.MustGet("login").(string)
	ctx.HTML(http.StatusOK, "showURL.html", gin.H{"login": login})
}

func (u *URL) GetUserURLsJSON(ctx *gin.Context) {
	login := auth.ValidLogin(ctx.MustGet("login").(string))
	urls, err := query.GetURLS(u.db, login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения ссылки пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, urls)
}

func (u *URL) CreateShortURLHandler(ctx *gin.Context) {

	loginValue, exists := ctx.Get("login")
	login := auth.ValidLogin(loginValue.(string))
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка получения данных"})
		return
	}
	login, ok := loginValue.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Некоректный логин"})
		return
	}

	// struct for information from url form
	var input struct {
		FullURL  string `json:"full_url"`
		ShortURL string `json:"short_url"` // необязательно
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		u.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка запроса"})
		return
	}

	//  validate url
	if !IsValidURL(input.FullURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некоректная ссылка"})
		return
	}

	//  generate short url if user didnt send
	if input.ShortURL == "" {
		input.ShortURL = ShortURL(8) // 8 — длина по умолчанию
	}

	account, err := query.GetAccount(u.db, login)
	if err != nil {
		u.logger.Error("Ошибка получения аккаунта", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}
	ok, err = query.IsExistsURL(u.db, account.Id, input.FullURL)
	if err != nil {
		u.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}
	if ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Такая ссылка уже существует"})
		return
	}
	// 6. save url in bd

	err = query.CreateURL(u.db, input.FullURL, input.ShortURL, account.Id)
	ok, Err := query.IsDuplicateKeyError(err)
	if ok {
		u.logger.Error(Err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Такая короткая ссылка уже существует ",
		})
		return
	}
	if err != nil {
		u.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса"})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Короткая ссылка создана!",
		"short_url": input.ShortURL,
	})
}

func (u *URL) FollowURLHandler(ctx *gin.Context) {
	short := ctx.Param("url")

	var url Table.URL
	var err error

	url.Full_url, err = query.GetURLByShortURL(u.db, short)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Страница не найдена"})
		return
	}

	ctx.Redirect(http.StatusFound, url.Full_url)
}

func (u *URL) DeleteURL(ctx *gin.Context) {
	short := ctx.Param("short_url")

	err := query.DeleteURLByShortURLL(u.db, short)
	if err != nil {
		u.logger.Info(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить ссылку. Попробуйте позже"})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Ссылка удалена успешна"})
	return
}

func IsValidURL(URL string) bool {
	check, err := url.Parse(URL)
	if err != nil || check.Host == "" || check.Scheme == "" {
		return false
	}

	ascii, err := idna.ToASCII(check.Host)
	if err != nil || ascii != check.Host {
		return false
	}

	if len(URL) > 2000 {
		return false
	}

	return check.Scheme == "https" || check.Scheme == "http"
}

func ShortURL(count int) string {
	short := make([]byte, count)
	for i := range short {
		short[i] = charset[SeedeRand.Intn(len(charset))]
	}
	return string(short)
}
