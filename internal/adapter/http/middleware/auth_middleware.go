package middleware

import (
	"net/http"
	"strings"

	"mecanica_xpto/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token ausente"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		c.Next()
	}
}
