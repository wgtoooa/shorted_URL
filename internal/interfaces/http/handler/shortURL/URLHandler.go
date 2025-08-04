package shortURL

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/idna"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"url_shortness/internal/domain/entities/Table"
	query2 "url_shortness/internal/infra/database/query"
	"url_shortness/internal/infra/product"
	"url_shortness/internal/interfaces/http/handler/auth"
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
	urls, err := query2.GetURLS(u.db, login)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Некорректный логин"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная ссылка"})
		return
	}

	//  generate short url if user didnt send
	if input.ShortURL == "" {
		input.ShortURL = ShortURL(8) // 8 — длина по умолчанию
	}

	account, err := product.GetUser(u.db, u.rdb, login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ok, err = query2.IsExistsURL(u.db, account.Id, input.FullURL)
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

	err = query2.CreateURL(u.db, input.FullURL, input.ShortURL, account.Id)
	ok, Err := query2.IsDuplicateKeyError(err)
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

	url.Full_url, err = query2.GetURLByShortURL(u.db, short)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Страница не найдена"})
		return
	}

	ctx.Redirect(http.StatusFound, url.Full_url)
}

func (u *URL) PatchURLHandler(ctx *gin.Context) {
	var request struct {
		OldShortURL string `json:"old_short_url" binding:"required"`
		NewShortURL string `json:"new_short_url" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		u.logger.Error("Ошибка валидации запроса", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат запроса",
			"details": err.Error(),
		})
		return
	}
	if request.OldShortURL == request.NewShortURL {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Новая ссылка должна отличаться от старой",
		})
		return
	}

	if err := query2.PutShortURL(u.db, request.OldShortURL, request.NewShortURL); err != nil {
		if ok, _ := query2.IsDuplicateKeyError(err); ok {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Такая короткая ссылка уже занята",
			})
			return
		}

		u.logger.Error("Ошибка при обновлении ссылки",
			zap.String("old", request.OldShortURL),
			zap.String("new", request.NewShortURL),
			zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Внутренняя ошибка сервера при изменении ссылки",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Короткая ссылка успешно изменена",
		"data": gin.H{
			"old_url": request.OldShortURL,
			"new_url": request.NewShortURL,
		},
	})
}

func (u *URL) DeleteURL(ctx *gin.Context) {
	short := ctx.Param("short_url")

	err := query2.DeleteURLByShortURLL(u.db, short)
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
