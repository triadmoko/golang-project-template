package usecase

import (
	"app/internal/features/user/delivery/http/dto"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/domain/entity"
	"app/internal/shared/domain/repository"
	"app/pkg"
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// UserUsecase defines the interface for user use cases
type UserUsecase interface {
	GetProfile(ctx context.Context, userID string) (*dto.UserResponse, int, error)
	UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*dto.UserResponse, int, error)
	GetUsers(ctx context.Context, queries map[string]string) ([]*dto.UserResponse, pkg.PaginationResponse, int, error)
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
func (u *userUsecase) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		u.logger.Error("u.userRepo.GetByID ", err)
		return nil, http.StatusNotFound, constants.GetError(constants.UserNotFound, lang)
	}

	// Convert to DTO response
	return dto.ToUserResponse(user), http.StatusOK, nil
}

// UpdateProfile updates user profile
func (u *userUsecase) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*dto.UserResponse, int, error) {
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

	// Build filter for update
	filter := entity.FilterUser{
		ID: userID,
	}

	// Save updated user
	if err := u.userRepo.Update(ctx, filter, user); err != nil {
		u.logger.Error("u.userRepo.Update ", err)
		return nil, http.StatusInternalServerError, constants.GetError(constants.FailedToUpdateUser, lang)
	}

	// Convert to DTO response
	return dto.ToUserResponse(user), http.StatusOK, nil
}

// GetUsers retrieves list of users with filtering and pagination
func (u *userUsecase) GetUsers(ctx context.Context, queries map[string]string) ([]*dto.UserResponse, pkg.PaginationResponse, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	// Build pagination
	pagination := pkg.PaginationBuilder(queries["per_page"], queries["page"])

	// Parse array filters
	var genders, roles []string
	if queries["genders"] != "" {
		genders = strings.Split(queries["genders"], ",")
	}
	if queries["roles"] != "" {
		roles = strings.Split(queries["roles"], ",")
	}

	// Build filter
	filter := entity.FilterUser{
		ID:        queries["id"],
		Email:     queries["email"],
		Username:  queries["username"],
		FirstName: queries["first_name"],
		LastName:  queries["last_name"],
		Status:    queries["status"],
		Gender:    queries["gender"],
		Role:      queries["role"],
		Provider:  queries["provider"],
		Genders:   genders,
		Roles:     roles,
		PerPage:   pagination.PerPage,
		Offset:    pagination.Offset,
	}

	// Get users from repository
	users, total, err := u.userRepo.List(ctx, filter)
	if err != nil {
		u.logger.Error("u.userRepo.List ", err)
		return nil, pkg.PaginationResponse{}, http.StatusInternalServerError, constants.GetError(constants.FailedToGetUsers, lang)
	}

	// Convert entity users to DTO response
	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, dto.ToUserResponse(user))
	}

	// Build pagination response
	totalPage := pkg.TotalPage(total, pagination.PerPage)
	paginationResponse := pkg.PaginationResponse{
		PerPage:   pagination.PerPage,
		TotalPage: totalPage,
		TotalData: total,
		Page:      pagination.Page,
	}

	return userResponses, paginationResponse, http.StatusOK, nil
}
