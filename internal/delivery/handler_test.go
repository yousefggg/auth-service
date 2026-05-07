package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yousefggg/auth-service/internal/usecase/mocks"
)

func TestHandler_Register(t *testing.T) {
	mockUC := new(mocks.AuthUseCase)
	handler := NewHandler(mockUC)

	t.Run("Success", func(t *testing.T) {
		userInput := map[string]string{
			"email":    "test@test.com",
			"password": "password123",
			"role":     "user",
		}
		body, _ := json.Marshal(userInput)

		mockUC.On("Register", mock.Anything, "test@test.com", "password123", "user").
			Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Invalid_JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte("invalid")))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}