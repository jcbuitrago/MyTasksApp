package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"mytasks/internal/auth"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")

		token, err := auth.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "claims inválidos"})
			return
		}

		sub, ok := claims["sub"].(float64) // JWT numérico llega como float64
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sub inválido"})
			return
		}
		c.Set("userID", int64(sub))
		c.Next()
	}
}

