package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization") //in json set header "Authorization" and here we him get
		var tokenString string

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Также пробуем взять токен из cookie для HTML-страниц
			cookie, err := ctx.Cookie("token")
			if err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			handleUnauthorized(ctx)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			handleUnauthorized(ctx)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handleUnauthorized(ctx)
			return
		}

		ctx.Set("login", claims["login"])
		ctx.Next()
	}
}

func handleUnauthorized(ctx *gin.Context) {
	accept := ctx.GetHeader("Accept")
	if strings.Contains(accept, "application/json") || ctx.Request.URL.Path == "/url/data" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	} else {
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
	}
}
