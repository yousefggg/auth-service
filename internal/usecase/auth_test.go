package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yousefggg/auth-service/internal/domain"
	"github.com/yousefggg/auth-service/internal/usecase/mocks"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthInteractor_Register(t *testing.T) {
	repo := new(mocks.UserRepository)
	tm := new(mocks.TokenManager)
	interactor := NewAuthInteractor(repo, tm)

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		role := domain.RoleUser

		repo.On("GetByEmail", mock.Anything, email).Return(nil, assert.AnError).Once()

		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == email && u.PasswordHash != password
		})).Return(nil).Once()

		err := interactor.Register(context.Background(), email, password, role)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Invalid_Input", func(t *testing.T) {
		err := interactor.Register(context.Background(), "bad", "12345678", domain.RoleUser)
		assert.Error(t, err)

		err = interactor.Register(context.Background(), "test@test.com", "123", domain.RoleUser)
		assert.Error(t, err)

		err = interactor.Register(context.Background(), "test@test.com", "12345678", "superman")
		assert.Error(t, err)
	})

	t.Run("User_Exists", func(t *testing.T) {
		email := "exists@test.com"
		existingUser := &domain.User{Email: email}

		repo.On("GetByEmail", mock.Anything, email).Return(existingUser, nil).Once()

		err := interactor.Register(context.Background(), email, "password", domain.RoleUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user already exists")
	})
}

func TestAuthInteractor_Login(t *testing.T) {
	repo := new(mocks.UserRepository)
	tm := new(mocks.TokenManager)
	interactor := NewAuthInteractor(repo, tm)

	t.Run("Success", func(t *testing.T) {
		email := "login@test.com"
		password := "correct_pass"
		userID := uuid.New()

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		existingUser := &domain.User{
			ID:           userID,
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "admin",
		}

		repo.On("GetByEmail", mock.Anything, email).Return(existingUser, nil).Once()

		tm.On("GenerateToken", userID, "admin").
			Return("valid_token", nil).Once()

		token, err := interactor.Login(context.Background(), email, password)

		assert.NoError(t, err)
		assert.Equal(t, "valid_token", token)
	})

	t.Run("Invalid_Credentials", func(t *testing.T) {
		repo.On("GetByEmail", mock.Anything, "wrong@test.com").
			Return(nil, assert.AnError).Once()

		token, err := interactor.Login(context.Background(), "wrong@test.com", "any")

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Wrong_Password", func(t *testing.T) {
		email := "pass@test.com"

		wrongHash, _ := bcrypt.GenerateFromPassword([]byte("real_password"), bcrypt.DefaultCost)
		user := &domain.User{
			Email:        email,
			PasswordHash: string(wrongHash),
		}

		repo.On("GetByEmail", mock.Anything, email).Return(user, nil).Once()

		token, err := interactor.Login(context.Background(), email, "wrong_password")

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Token_Generation_Fail", func(t *testing.T) {
		email := "jwt@test.com"
		userID := uuid.New()

		hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		user := &domain.User{
			ID:           userID,
			Email:        email,
			PasswordHash: string(hash),
			Role:         "user",
		}

		repo.On("GetByEmail", mock.Anything, email).Return(user, nil).Once()

		tm.On("GenerateToken", userID, "user").
			Return("", assert.AnError).Once()

		token, err := interactor.Login(context.Background(), email, "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}