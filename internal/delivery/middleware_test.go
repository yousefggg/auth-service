package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем его в наш Middleware
	middleware := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Проверяем, что запрос прошел сквозь middleware и дошел до хендлера
	assert.Equal(t, http.StatusOK, rec.Code)
}