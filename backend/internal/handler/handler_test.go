package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newHandlerTestRouter() *gin.Engine {
	gin.SetMode(
		gin.TestMode,
	)

	router :=
		gin.New()

	router.GET(
		"/api/v1/health",
		Health,
	)

	router.GET(
		"/api/v1/check",
		Check,
	)

	router.GET(
		"/api/v1/check-multiple",
		CheckMultiple,
	)

	return router
}

func TestHealth(
	t *testing.T,
) {
	router :=
		newHandlerTestRouter()

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/api/v1/health",
			nil,
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
			"status = %d, want 200",
			recorder.Code,
		)
	}

	var response map[string]any

	if err :=
		json.Unmarshal(
			recorder.Body.Bytes(),
			&response,
		); err != nil {
		t.Fatalf(
			"JSON tidak valid: %v",
			err,
		)
	}

	if response["status"] !=
		"ok" {
		t.Fatalf(
			"status response = %v, want ok",
			response["status"],
		)
	}

	if response["service"] !=
		"gowatch-api" {
		t.Fatalf(
			"service = %v, want gowatch-api",
			response["service"],
		)
	}
}

func TestCheck_MissingURL(
	t *testing.T,
) {
	router :=
		newHandlerTestRouter()

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/api/v1/check",
			nil,
		)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want 400",
			recorder.Code,
		)
	}
}

func TestCheck_BlocksLocalhost(
	t *testing.T,
) {
	router :=
		newHandlerTestRouter()

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/api/v1/check?url=http%3A%2F%2F127.0.0.1%3A8081",
			nil,
		)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want 400. body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !strings.Contains(
		strings.ToLower(
			recorder.Body.String(),
		),
		"tidak diizinkan",
	) {
		t.Fatalf(
			"response tidak menunjukkan URL diblokir: %s",
			recorder.Body.String(),
		)
	}
}

func TestCheckMultiple_Empty(
	t *testing.T,
) {
	router :=
		newHandlerTestRouter()

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/api/v1/check-multiple",
			nil,
		)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want 400",
			recorder.Code,
		)
	}
}

func TestCheckMultiple_MaximumTen(
	t *testing.T,
) {
	router :=
		newHandlerTestRouter()

	query :=
		"/api/v1/check-multiple?"

	for index :=
		0; index < 11; index++ {
		if index > 0 {
			query += "&"
		}

		query +=
			"url=https%3A%2F%2F1.1.1.1"
	}

	request :=
		httptest.NewRequest(
			http.MethodGet,
			query,
			nil,
		)

	recorder :=
		httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		request,
	)

	if recorder.Code !=
		http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want 400. body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !strings.Contains(
		recorder.Body.String(),
		"maksimal 10",
	) {
		t.Fatalf(
			"pesan limit tidak ditemukan: %s",
			recorder.Body.String(),
		)
	}
}
