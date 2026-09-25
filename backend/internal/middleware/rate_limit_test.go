package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimit(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	router :=
		gin.New()

	if err :=
		router.SetTrustedProxies(
			nil,
		); err != nil {
		t.Fatalf(
			"SetTrustedProxies gagal: %v",
			err,
		)
	}

	router.GET(
		"/limited",
		RateLimit(
			2,
			time.Minute,
		),
		func(
			c *gin.Context,
		) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"status": "ok",
				},
			)
		},
	)

	for index :=
		1; index <= 2; index++ {
		request :=
			httptest.NewRequest(
				http.MethodGet,
				"/limited",
				nil,
			)

		request.RemoteAddr =
			"8.8.8.8:12345"

		recorder :=
			httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			request,
		)

		if recorder.Code !=
			http.StatusOK {
			t.Fatalf(
				"request %d: status = %d, want 200",
				index,
				recorder.Code,
			)
		}
	}

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/limited",
			nil,
		)

	request.RemoteAddr =
		"8.8.8.8:12345"

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusTooManyRequests {
		t.Fatalf(
			"status = %d, want 429",
			recorder.Code,
		)
	}

	if recorder.Header().
		Get(
			"Retry-After",
		) == "" {
		t.Fatal(
			"Retry-After header seharusnya tersedia",
		)
	}
}

func TestRateLimit_DifferentIPs(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	router :=
		gin.New()

	if err :=
		router.SetTrustedProxies(
			nil,
		); err != nil {
		t.Fatalf(
			"SetTrustedProxies gagal: %v",
			err,
		)
	}

	router.GET(
		"/limited",
		RateLimit(
			1,
			time.Minute,
		),
		func(
			c *gin.Context,
		) {
			c.Status(
				http.StatusOK,
			)
		},
	)

	addresses :=
		[]string{
			"8.8.8.8:1000",
			"1.1.1.1:1000",
		}

	for _, address := range addresses {
		request :=
			httptest.NewRequest(
				http.MethodGet,
				"/limited",
				nil,
			)

		request.RemoteAddr =
			address

		recorder :=
			httptest.NewRecorder()

		router.ServeHTTP(
			recorder,
			request,
		)

		if recorder.Code !=
			http.StatusOK {
			t.Fatalf(
				"IP %s mendapat status %d, want 200",
				address,
				recorder.Code,
			)
		}
	}
}
