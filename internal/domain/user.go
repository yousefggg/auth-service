package domain

import (
	"context"
	"time"
	"github.com/google/uuid"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           uuid.UUID    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` 
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthUsecase interface {
	Register(ctx context.Context, email, password, role string) error
	Login(ctx context.Context, email, password string) (string, error) 
}
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}