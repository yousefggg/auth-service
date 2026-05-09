package delivery

import (
	"encoding/json"
	"net/http"
	"context"
	"github.com/yousefggg/common-lib/pkg/errors"
	"github.com/yousefggg/common-lib/pkg/logger"
	"github.com/yousefggg/common-lib/pkg/dto"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password, role string) error
	Login(ctx context.Context, email, password string) (string, error)
	ParseToken(token string) (*dto.Claims, error)
}

type Handler struct {
	authUseCase AuthUseCase 
}

func NewHandler(authUseCase AuthUseCase) *Handler {
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
// Register godoc
// @Summary      Регистрация
// @Description  Создает нового пользователя в системе
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body      domain.User true "Данные пользователя (email, password, role)"
// @Success      201  {object}  map[string]string "message: user created"
// @Failure      400  {object}  map[string]string "error: invalid input"
// @Failure      500  {object}  map[string]string "error: internal server error"
// @Router       /auth/register [post]
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
// Login godoc
// @Summary      Вход
// @Description  Проверяет данные и возвращает JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body      domain.User true "Данные для входа (email, password)"
// @Success      200  {object}  map[string]string "token: JWT_TOKEN_HERE"
// @Failure      401  {object}  map[string]string "error: unauthorized"
// @Router       /auth/login [post]
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
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")

    claims, err := h.authUseCase.ParseToken(token)
    if err != nil {
        json.NewEncoder(w).Encode(dto.ValidateTokenResponse{Valid: false})
        return
    }
    resp := dto.ValidateTokenResponse{
        UserID: claims.UserID,
        Role:   claims.Role,
        Valid:  true,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}