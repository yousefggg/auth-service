package usecase

import (
    "context"
    "strings"
    "github.com/google/uuid"
    "github.com/yousefggg/auth-service/internal/domain"
    "github.com/yousefggg/common-lib/pkg/errors"
    "github.com/yousefggg/common-lib/pkg/logger"
    "golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
    GenerateToken(userID uuid.UUID, role string) (string, error)
}
type AuthInteractor struct {
    repo         domain.UserRepository
    tokenManager TokenManager
}

func NewAuthInteractor(repo domain.UserRepository, tm TokenManager) *AuthInteractor {
    return &AuthInteractor{
        repo:         repo,
        tokenManager: tm,
    }
}

func (a *AuthInteractor) Register(ctx context.Context, email, password, role string) error {
	if !strings.Contains(email, "@") || len(email) < 5 {
		logger.Warn("Registration failed: invalid email format", "email", email)
		return errors.NewErr("AUTH_INVALID_EMAIL", "invalid email format", nil)
	}

	if len(password) < 8 {
		logger.Warn("Registration failed: password too short", "email", email)
		return errors.NewErr("AUTH_SHORT_PASSWORD", "password must be at least 8 characters", nil)
	}

	if role != domain.RoleUser && role != domain.RoleAdmin {
		logger.Warn("Registration failed: invalid role", "role", role)
		return errors.NewErr("AUTH_INVALID_ROLE", "unsupported role", nil)
	}

	existingUser, _ := a.repo.GetByEmail(ctx, email)
	if existingUser != nil {
		logger.Warn("Registration failed: user already exists", "email", email)
		return errors.NewErr("AUTH_USER_EXISTS", "user already exists", nil)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", "error", err)
		return errors.NewErr("INTERNAL_ERROR", "failed to process security data", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := a.repo.Create(ctx, user); err != nil {
		logger.Error("Failed to create user", "email", email, "error", err)
		return err
	}

	logger.Info("User registered successfully", "email", email, "role", role)
	return nil
}

func (a *AuthInteractor) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		logger.Warn("Login failed: user not found", "email", email)
		return "", errors.NewErr("AUTH_FAILED", "invalid credentials", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logger.Warn("Login failed: wrong password", "email", email)
		return "", errors.NewErr("AUTH_FAILED", "invalid credentials", nil)
	}

	token, err := a.tokenManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		logger.Error("Failed to generate JWT", "user_id", user.ID, "error", err)
		return "", errors.NewErr("INTERNAL_ERROR", "failed to create session", err)
	}

	logger.Info("User logged in successfully", "user_id", user.ID, "email", email)
	return token, nil
}