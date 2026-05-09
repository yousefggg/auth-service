package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yousefggg/auth-service/internal/domain"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type TokenManager struct {
	mock.Mock
}

func (m *TokenManager) GenerateToken(userID uuid.UUID, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *TokenManager) ValidateToken(token string) (uuid.UUID, string, error) {
	args := m.Called(token)

	var userID uuid.UUID
	if args.Get(0) != nil {
		userID = args.Get(0).(uuid.UUID)
	}

	return userID, args.String(1), args.Error(2)
}