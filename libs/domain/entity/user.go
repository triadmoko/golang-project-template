package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user entity in the domain layer
type User struct {
	ID        string         `json:"id" gorm:"type:varchar(36);primaryKey"`
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

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user entity with generated UUID
func NewUser(email, username, password, firstName, lastName string) *User {
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}
}

// BeforeCreate hook to ensure UUID is set
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
