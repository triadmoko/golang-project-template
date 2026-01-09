package usecase

import (
	"app/internal/features/auth/delivery/http/dto"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/domain/entity"
	"app/internal/shared/domain/repository"
	"app/pkg/crypto"
	"app/pkg/jwt"
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

// AuthUsecase defines the interface for authentication use cases
type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, int, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, int, error)
}

// authUsecase implements AuthUsecase interface
type authUsecase struct {
	userRepo repository.UserRepository
	logger   *logrus.Logger
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo repository.UserRepository, logger *logrus.Logger) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Register creates a new user
func (a *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	// Check if user already exists by email
	existingUser, _ := a.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		a.logger.Error("a.userRepo.GetByEmail: user already exists")
		return nil, http.StatusBadRequest, constants.GetError(constants.UserAlreadyExists, lang)
	}

	// Check if username is taken
	existingUser, _ = a.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		a.logger.Error("a.userRepo.GetByUsername: username already taken")
		return nil, http.StatusBadRequest, constants.GetError(constants.UsernameAlreadyTaken, lang)
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		a.logger.Error("crypto.HashPassword ", err)
		return nil, http.StatusInternalServerError, constants.GetError(constants.FailedToHashPassword, lang)
	}

	// Create user entity using shared entity
	user := entity.NewUser(req.Email, req.Username, hashedPassword, req.FirstName, req.LastName)

	// Save user
	if err := a.userRepo.Create(ctx, user); err != nil {
		a.logger.Error("a.userRepo.Create ", err)
		return nil, http.StatusInternalServerError, constants.GetError(constants.FailedToCreateUser, lang)
	}

	// Convert to DTO response
	return dto.ToRegisterResponse(user), http.StatusCreated, nil
}

// Login authenticates a user
func (a *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	// Get user by email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		a.logger.Error("a.userRepo.GetByEmail ", err)
		return nil, http.StatusUnauthorized, constants.GetError(constants.InvalidCredentials, lang)
	}

	// Verify password
	if err := crypto.VerifyPassword(user.Password, req.Password); err != nil {
		a.logger.Error("crypto.VerifyPassword ", err)
		return nil, http.StatusUnauthorized, constants.GetError(constants.InvalidCredentials, lang)
	}

	// Generate token with string UUID
	token, err := jwt.GenerateToken(jwt.UserPayload{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		a.logger.Error("jwt.GenerateToken ", err)
		return nil, http.StatusInternalServerError, constants.GetError(constants.FailedToGenerateToken, lang)
	}

	// Convert to DTO response
	return &dto.LoginResponse{
		User:  dto.ToRegisterResponse(user),
		Token: token,
	}, http.StatusOK, nil
}
