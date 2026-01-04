package usecase

import (
	"app/internal/features/user/delivery/http/dto"
	"app/internal/shared/domain/entity"
	"app/internal/shared/domain/repository"
	"context"
	"net/http"

	domainError "app/internal/shared/domain/error"

	"github.com/sirupsen/logrus"
)

// UserUsecase defines the interface for user use cases
type UserUsecase interface {
	GetProfile(ctx context.Context, userID string) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*entity.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo repository.UserRepository
	logger   *logrus.Logger
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository, logger *logrus.Logger) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetProfile retrieves user profile
func (u *userUsecase) GetProfile(ctx context.Context, userID string) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("u.userRepo.GetByID ", err)
		return nil, domainError.NewCustomError(http.StatusNotFound, "user not found", domainError.ErrUserNotFound)
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// UpdateProfile updates user profile
func (u *userUsecase) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("u.userRepo.GetByID ", err)
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
		u.logger.Error("u.userRepo.Update ", err)
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
		u.logger.Error("u.userRepo.List ", err)
		return nil, domainError.NewCustomError(http.StatusInternalServerError, "failed to get users", err)
	}

	// Remove passwords from response
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}
