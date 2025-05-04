package shortURL

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/idna"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"url_shortness/internal/repository/Table"
	"url_shortness/internal/repository/database"
	"url_shortness/internal/repository/database/query"
)

const (
	lenShortURL = 7
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var SeedeRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func ShowURLhandler(ctx *gin.Context) {
	login := ctx.MustGet("login").(string)
	ctx.HTML(http.StatusOK, "showURL.html", gin.H{"login": login})
}
func GetUserURLsJSON(ctx *gin.Context) {
	login := ctx.MustGet("login").(string)
	urls, err := query.GetURLS(database.DB, login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed get user url"})
		return
	}

	ctx.JSON(http.StatusOK, urls)
}

func CreateShortURLhandler(ctx *gin.Context) {
	// 1. Получаем логин пользователя из контекста (установлен в middleware)
	loginValue, exists := ctx.Get("login")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	login, ok := loginValue.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid login value"})
		return
	}

	// 2. Структура для привязки данных из запроса
	var input struct {
		FullURL     string `json:"full_url"`
		ShortURL    string `json:"short_url"`   // необязательно
		Description string `json:"description"` // необязательно
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 3. Валидация URL
	if !IsValidURL(input.FullURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "url not valid"})
		return
	}

	// 4. Генерация короткого URL, если он не передан
	if input.ShortURL == "" {
		input.ShortURL = ShortURL(8) // 8 — длина по умолчанию
	}

	// 5. Получаем account_id по логину
	account, err := query.GetAccount(database.DB, login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	// 6. Сохраняем ссылку в БД
	err = query.CreateURL(database.DB, input.FullURL, input.ShortURL, account.Id, input.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create url"})
		return
	}

	// 7. Ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message":   "short URL created",
		"short_url": input.ShortURL,
	})
}

func FolowURLHandler(ctx *gin.Context) {
	short := ctx.Param("url")

	var url Table.URL
	var err error

	url.Full_url, err = query.GetURLByShortURL(database.DB, short)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found page"})
		return
	}

	ctx.Redirect(http.StatusFound, url.Full_url)
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
