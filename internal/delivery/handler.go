package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/yousefggg/auth-service/internal/usecase"
	"github.com/yousefggg/common-lib/pkg/errors"
	"github.com/yousefggg/common-lib/pkg/logger"
)

type Handler struct {
	authUseCase *usecase.AuthInteractor
}

func NewHandler(authUseCase *usecase.AuthInteractor) *Handler {
	return &Handler{
		authUseCase: authUseCase,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode register request", "error", err)
		h.sendError(w, errors.NewErr("BAD_REQUEST", "invalid request body", err))
		return
	}

	err := h.authUseCase.Register(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		h.sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode login request", "error", err)
		h.sendError(w, errors.NewErr("BAD_REQUEST", "invalid request body", err))
		return
	}

	token, err := h.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) sendError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	if customErr, ok := err.(*errors.CustomError); ok {
		switch customErr.Code {
		case "AUTH_INVALID_EMAIL", "AUTH_SHORT_PASSWORD", "AUTH_INVALID_ROLE", "BAD_REQUEST":
			status = http.StatusBadRequest
		case "AUTH_USER_EXISTS":
			status = http.StatusConflict
		case "AUTH_FAILED":
			status = http.StatusUnauthorized
		}

		h.sendJSON(w, status, map[string]string{
			"code":    customErr.Code,
			"message": customErr.Message,
		})
		return
	}

	h.sendJSON(w, status, map[string]string{"message": "internal server error"})
}