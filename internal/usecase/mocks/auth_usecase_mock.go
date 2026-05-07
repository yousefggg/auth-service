package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type AuthUseCase struct {
	mock.Mock
}

func (m *AuthUseCase) Register(ctx context.Context, email, password, role string) error {
	args := m.Called(ctx, email, password, role)
	return args.Error(0)
}

func (m *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}