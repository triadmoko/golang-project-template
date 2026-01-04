package usecase

import (
	"app/internal/features/user/delivery/http/dto"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/domain/entity"
	"app/internal/shared/domain/repository"
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

// UserUsecase defines the interface for user use cases
type UserUsecase interface {
	GetProfile(ctx context.Context, userID string) (*entity.User, int, error)
	UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*entity.User, int, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*entity.User, int, error)
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
func (u *userUsecase) GetProfile(ctx context.Context, userID string) (*entity.User, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("u.userRepo.GetByID ", err)
		return nil, http.StatusNotFound, constants.GetError(constants.UserNotFound, lang)
	}

	// Remove password from response
	user.Password = ""
	return user, http.StatusOK, nil
}

// UpdateProfile updates user profile
func (u *userUsecase) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*entity.User, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("u.userRepo.GetByID ", err)
		return nil, http.StatusNotFound, constants.GetError(constants.UserNotFound, lang)
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
		return nil, http.StatusInternalServerError, constants.GetError(constants.FailedToUpdateUser, lang)
	}

	// Remove password from response
	user.Password = ""
	return user, http.StatusOK, nil
}

// GetUsers retrieves list of users
func (u *userUsecase) GetUsers(ctx context.Context, limit, offset int) ([]*entity.User, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		u.logger.Error("u.userRepo.List ", err)
		return nil, http.StatusInternalServerError, constants.GetError(constants.FailedToGetUsers, lang)
	}

	// Remove passwords from response
	for _, user := range users {
		user.Password = ""
	}

	return users, http.StatusOK, nil
}
