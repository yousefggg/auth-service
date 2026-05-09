package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yousefggg/common-lib/pkg/dto"
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

func (m *AuthUseCase) ParseToken(token string) (*dto.Claims, error) {
	args := m.Called(token)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.Claims), args.Error(1)
}