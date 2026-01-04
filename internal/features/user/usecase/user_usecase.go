package usecase

import (
	"app/internal/shared/domain/entity"
	"app/internal/shared/domain/repository"
	"context"
	"net/http"

	domainError "app/internal/shared/domain/error"
)

// UserUsecase defines the interface for user use cases
type UserUsecase interface {
	GetProfile(ctx context.Context, userID string) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*entity.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// UpdateProfileRequest represents the request for updating user profile
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetProfile retrieves user profile
func (u *userUsecase) GetProfile(ctx context.Context, userID string) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domainError.NewCustomError(http.StatusNotFound, "user not found", domainError.ErrUserNotFound)
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// UpdateProfile updates user profile
func (u *userUsecase) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domainError.NewCustomError(http.StatusNotFound, "user not found", domainError.ErrUserNotFound)
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	// Save updated user
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, domainError.NewCustomError(http.StatusInternalServerError, "failed to update user", err)
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// GetUsers retrieves list of users
func (u *userUsecase) GetUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, domainError.NewCustomError(http.StatusInternalServerError, "failed to get users", err)
	}

	// Remove passwords from response
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}
