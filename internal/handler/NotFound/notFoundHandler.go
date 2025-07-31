package NotFound

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFoundHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "404.html", nil)
}
