package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "hello.html", nil)
}
