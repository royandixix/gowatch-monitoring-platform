package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/royandi/gowatch/backend/internal/auth"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		if authorizationHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization token diperlukan",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(
			authorizationHeader,
			" ",
			2,
		)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "format authorization tidak valid",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token tidak boleh kosong",
			})
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(
			tokenString,
			jwtSecret,
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token tidak valid atau sudah kedaluwarsa",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
