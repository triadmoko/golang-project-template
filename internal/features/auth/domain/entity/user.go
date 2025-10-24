package entity

import (
	"app/internal/features/auth/delivery/http/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user entity in the domain layer
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID      string         `json:"uuid" gorm:"type:varchar(36);uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	FirstName string         `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName  string         `json:"last_name" gorm:"type:varchar(100);not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// NewUser creates a new user entity
func NewUser(req dto.RegisterRequest) *User {
	return &User{
		UUID:      uuid.New().String(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
