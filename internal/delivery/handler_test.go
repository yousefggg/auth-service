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
		bodyData := RegisterRequest{
			Email:    "test@test.com",
			Password: "password123",
			Role:     "user",
		}

		body, err := json.Marshal(bodyData)
		assert.NoError(t, err)

		mockUC.
			On("Register",
				mock.Anything,
				bodyData.Email,
				bodyData.Password,
				bodyData.Role,
			).
			Return(nil).
			Once()

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Invalid_JSON", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodPost,
			"/auth/register",
			bytes.NewBuffer([]byte("{invalid-json}")),
		)

		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}