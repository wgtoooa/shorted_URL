package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) HelloHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "hello.html", nil)
}
