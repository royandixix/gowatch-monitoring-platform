package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	Count       int
	WindowStart time.Time
}

func RateLimit(
	maxRequests int,
	window time.Duration,
) gin.HandlerFunc {
	if maxRequests < 1 {
		panic(
			"maxRequests harus lebih besar dari 0",
		)
	}

	if window <= 0 {
		panic(
			"window rate limit harus lebih besar dari 0",
		)
	}

	var mutex sync.Mutex

	clients :=
		make(
			map[string]*rateLimitEntry,
		)

	lastCleanup :=
		time.Now()

	return func(
		c *gin.Context,
	) {
		now :=
			time.Now()

		clientIP :=
			c.ClientIP()

		mutex.Lock()

		if now.Sub(
			lastCleanup,
		) >= window {
			for key, entry := range clients {
				if now.Sub(
					entry.WindowStart,
				) >= window {
					delete(
						clients,
						key,
					)
				}
			}

			lastCleanup =
				now
		}

		entry,
			exists :=
			clients[clientIP]

		if !exists ||
			now.Sub(
				entry.WindowStart,
			) >= window {
			entry =
				&rateLimitEntry{
					Count: 0,

					WindowStart: now,
				}

			clients[clientIP] =
				entry
		}

		if entry.Count >=
			maxRequests {
			retryAfter :=
				int(
					math.Ceil(
						entry.WindowStart.
							Add(window).
							Sub(now).
							Seconds(),
					),
				)

			if retryAfter < 1 {
				retryAfter = 1
			}

			mutex.Unlock()

			c.Header(
				"Retry-After",
				strconv.Itoa(
					retryAfter,
				),
			)

			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				gin.H{
					"error": "terlalu banyak request, coba lagi beberapa saat",
				},
			)

			return
		}

		entry.Count++

		remaining :=
			maxRequests -
				entry.Count

		resetAt :=
			entry.WindowStart.
				Add(window).
				Unix()

		mutex.Unlock()

		c.Header(
			"X-RateLimit-Limit",
			strconv.Itoa(
				maxRequests,
			),
		)

		c.Header(
			"X-RateLimit-Remaining",
			strconv.Itoa(
				remaining,
			),
		)

		c.Header(
			"X-RateLimit-Reset",
			strconv.FormatInt(
				resetAt,
				10,
			),
		)

		c.Next()
	}
}
