package service

import (
	"app/internal/features/auth/domain/entity"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	GenerateToken(user *entity.User) (string, error)
	ValidateToken(token string) (*entity.User, error)
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
