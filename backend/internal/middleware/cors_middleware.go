package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(
	allowedOrigin string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin == allowedOrigin {
			c.Header(
				"Access-Control-Allow-Origin",
				allowedOrigin,
			)

			c.Header(
				"Access-Control-Allow-Credentials",
				"true",
			)

			c.Header(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)

			c.Header(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)

			c.Header(
				"Vary",
				"Origin",
			)
		}

		if c.Request.Method == http.MethodOptions {
			if origin != allowedOrigin {
				c.AbortWithStatus(
					http.StatusForbidden,
				)

				return
			}

			c.AbortWithStatus(
				http.StatusNoContent,
			)

			return
		}

		c.Next()
	}
}
