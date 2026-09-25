package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/royandi/gowatch/backend/internal/model"
	"github.com/royandi/gowatch/backend/internal/monitor"
)

func Home(
	c *gin.Context,
) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "GoWatch API berjalan 🚀",

			"version": "1.0.0",

			"endpoints": []string{
				"/api/v1/health",
				"/api/v1/check?url=https://google.com",
				"/api/v1/check-multiple?url=https://google.com&url=https://github.com",
			},
		},
	)
}

func Health(
	c *gin.Context,
) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",

			"service": "gowatch-api",
		},
	)
}

func Check(
	c *gin.Context,
) {
	targetURL :=
		c.Query(
			"url",
		)

	if targetURL == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "parameter url wajib diisi",
			},
		)

		return
	}

	normalizedURL, err :=
		monitor.NormalizeAndValidateURL(
			c.Request.Context(),
			targetURL,
		)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "URL tidak diizinkan: " +
					err.Error(),
			},
		)

		return
	}

	result :=
		monitor.CheckWebsiteContext(
			c.Request.Context(),
			normalizedURL,
		)

	c.JSON(
		http.StatusOK,
		result,
	)
}

func CheckMultiple(
	c *gin.Context,
) {
	targetURLs :=
		c.QueryArray(
			"url",
		)

	if len(targetURLs) == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "minimal satu URL wajib diisi",
			},
		)

		return
	}

	const maxURLs = 10

	if len(targetURLs) >
		maxURLs {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "maksimal 10 URL dalam satu request",
			},
		)

		return
	}

	normalizedURLs :=
		make(
			[]string,
			0,
			len(targetURLs),
		)

	for _, targetURL := range targetURLs {
		normalizedURL, err :=
			monitor.NormalizeAndValidateURL(
				c.Request.Context(),
				targetURL,
			)

		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "URL tidak diizinkan: " +
						targetURL +
						" - " +
						err.Error(),
				},
			)

			return
		}

		normalizedURLs = append(
			normalizedURLs,
			normalizedURL,
		)
	}

	start :=
		time.Now()

	results :=
		monitor.CheckWebsitesConcurrentlyContext(
			c.Request.Context(),
			normalizedURLs,
		)

	response :=
		model.MultipleCheckResponse{
			Total: len(results),

			TotalDurationMs: time.Since(
				start,
			).Milliseconds(),

			Results: results,
		}

	c.JSON(
		http.StatusOK,
		response,
	)
}
