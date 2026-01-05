package usecase

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"monorepo/libs/crypto"
	"monorepo/libs/domain/entity"
	"monorepo/libs/domain/repository"
	"monorepo/libs/errors"
	"monorepo/libs/httputil/middleware"
	"monorepo/libs/jwt"
	"monorepo/services/api-gateway/internal/features/auth/delivery/http/dto"
)

// AuthUsecase defines the interface for authentication use cases
type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, int, error)
	Login(ctx context.Context, req dto.LoginRequest) (*LoginResponse, int, error)
}

// authUsecase implements AuthUsecase interface
type authUsecase struct {
	userRepo  repository.UserRepository
	logger    *logrus.Logger
	jwtSecret string
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo repository.UserRepository, logger *logrus.Logger, jwtSecret string) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	User  *entity.User `json:"user"`
	Token string       `json:"token"`
}

// Register creates a new user
func (a *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	// Check if user already exists by email
	existingUser, _ := a.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		a.logger.Error("a.userRepo.GetByEmail: user already exists")
		return nil, http.StatusBadRequest, errors.GetError(errors.UserAlreadyExists, lang)
	}

	// Check if username is taken
	existingUser, _ = a.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		a.logger.Error("a.userRepo.GetByUsername: username already taken")
		return nil, http.StatusBadRequest, errors.GetError(errors.UsernameAlreadyTaken, lang)
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		a.logger.Error("crypto.HashPassword ", err)
		return nil, http.StatusInternalServerError, errors.GetError(errors.FailedToHashPassword, lang)
	}

	// Create user entity using shared entity
	user := entity.NewUser(req.Email, req.Username, hashedPassword, req.FirstName, req.LastName)

	// Save user
	if err := a.userRepo.Create(ctx, user); err != nil {
		a.logger.Error("a.userRepo.Create ", err)
		return nil, http.StatusInternalServerError, errors.GetError(errors.FailedToCreateUser, lang)
	}

	// Remove password from response
	user.Password = ""
	return user, http.StatusCreated, nil
}

// Login authenticates a user
func (a *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*LoginResponse, int, error) {
	lang := middleware.GetLangFromContext(ctx)

	// Get user by email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		a.logger.Error("a.userRepo.GetByEmail ", err)
		return nil, http.StatusUnauthorized, errors.GetError(errors.InvalidCredentials, lang)
	}

	// Verify password
	if err := crypto.VerifyPassword(user.Password, req.Password); err != nil {
		a.logger.Error("crypto.VerifyPassword ", err)
		return nil, http.StatusUnauthorized, errors.GetError(errors.InvalidCredentials, lang)
	}

	// Generate token with string UUID
	token, err := jwt.GenerateToken(a.jwtSecret, jwt.UserPayload{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		a.logger.Error("jwt.GenerateToken ", err)
		return nil, http.StatusInternalServerError, errors.GetError(errors.FailedToGenerateToken, lang)
	}

	// Remove password from response
	user.Password = ""

	return &LoginResponse{
		User:  user,
		Token: token,
	}, http.StatusOK, nil
}
