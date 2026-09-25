package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/royandi/gowatch/backend/internal/auth"
)

func createProtectedTestRouter(
	secret string,
) *gin.Engine {
	gin.SetMode(
		gin.TestMode,
	)

	router :=
		gin.New()

	router.GET(
		"/protected",
		AuthMiddleware(
			secret,
		),
		func(
			c *gin.Context,
		) {
			userID, _ :=
				c.Get(
					"user_id",
				)

			email, _ :=
				c.Get(
					"user_email",
				)

			c.JSON(
				http.StatusOK,
				gin.H{
					"user_id": userID,

					"email": email,
				},
			)
		},
	)

	return router
}

func TestAuthMiddleware_NoToken(
	t *testing.T,
) {
	router :=
		createProtectedTestRouter(
			"test-secret",
		)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusUnauthorized {
		t.Fatalf(
			"status = %d, want 401",
			recorder.Code,
		)
	}
}

func TestAuthMiddleware_InvalidFormat(
	t *testing.T,
) {
	router :=
		createProtectedTestRouter(
			"test-secret",
		)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

	request.Header.Set(
		"Authorization",
		"invalid-token",
	)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusUnauthorized {
		t.Fatalf(
			"status = %d, want 401",
			recorder.Code,
		)
	}
}

func TestAuthMiddleware_ValidToken(
	t *testing.T,
) {
	const secret = "test-secret"

	token, err :=
		auth.GenerateToken(
			99,
			"tester@example.com",
			secret,
		)

	if err != nil {
		t.Fatalf(
			"GenerateToken gagal: %v",
			err,
		)
	}

	router :=
		createProtectedTestRouter(
			secret,
		)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

	request.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusOK {
		t.Fatalf(
			"status = %d, want 200. body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestAuthMiddleware_WrongSecret(
	t *testing.T,
) {
	token, err :=
		auth.GenerateToken(
			1,
			"tester@example.com",
			"secret-one",
		)

	if err != nil {
		t.Fatalf(
			"GenerateToken gagal: %v",
			err,
		)
	}

	router :=
		createProtectedTestRouter(
			"secret-two",
		)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

	request.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusUnauthorized {
		t.Fatalf(
			"status = %d, want 401",
			recorder.Code,
		)
	}
}
