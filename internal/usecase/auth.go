package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yousefggg/auth-service/internal/domain"
	"github.com/yousefggg/common-lib/pkg/jwt"
	"github.com/yousefggg/common-lib/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor struct {
	repo         domain.UserRepository
	tokenManager *jwt.TokenManager
}

func NewAuthInteractor(repo domain.UserRepository, tm *jwt.TokenManager) *AuthInteractor {
	return &AuthInteractor{
		repo:         repo,
		tokenManager: tm,
	}
}

func (a *AuthInteractor) Register(ctx context.Context, email, password, role string) error {
	if !strings.Contains(email, "@") || len(email) < 5 {
		logger.Warn("Registration failed: invalid email format", "email", email)
		return errors.New("invalid email format")
	}

	if len(password) < 8 {
		logger.Warn("Registration failed: password too short", "email", email)
		return errors.New("password must be at least 8 characters long")
	}

	if role != "user" && role != "admin" {
		logger.Warn("Registration failed: invalid role", "role", role)
		return errors.New("invalid role")
	}

	existingUser, _ := a.repo.GetByEmail(ctx, email)
	if existingUser != nil {
		logger.Warn("Registration failed: user already exists", "email", email)
		return errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password during registration", "error", err)
		return fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := a.repo.Create(ctx, user); err != nil {
		logger.Error("Failed to create user in database", "email", email, "error", err)
		return err
	}

	logger.Info("User registered successfully", "email", email, "role", role)
	return nil
}

func (a *AuthInteractor) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		logger.Warn("Login failed: user not found", "email", email)
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logger.Warn("Login failed: wrong password", "email", email)
		return "", errors.New("invalid credentials")
	}

	token, err := a.tokenManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		logger.Error("Failed to generate JWT token", "user_id", user.ID, "error", err)
		return "", fmt.Errorf("generate token: %w", err)
	}

	logger.Info("User logged in successfully", "user_id", user.ID, "email", email)
	return token, nil
}