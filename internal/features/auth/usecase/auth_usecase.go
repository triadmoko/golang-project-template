package usecase

import (
	"app/internal/features/auth/delivery/http/dto"
	"app/internal/shared/domain/entity"
	"app/internal/shared/domain/repository"
	"app/pkg/crypto"
	"app/pkg/jwt"
	"context"
	"net/http"

	domainError "app/internal/shared/domain/error"
)

// AuthUsecase defines the interface for authentication use cases
type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*LoginResponse, error)
}

// authUsecase implements AuthUsecase interface
type authUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	User  *entity.User `json:"user"`
	Token string       `json:"token"`
}

// Register creates a new user
func (a *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, error) {
	// Check if user already exists by email
	existingUser, _ := a.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domainError.NewCustomError(http.StatusBadRequest, "user already exists", domainError.ErrUserAlreadyExists)
	}

	// Check if username is taken
	existingUser, _ = a.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, domainError.NewCustomError(http.StatusBadRequest, "username already taken", domainError.ErrUserAlreadyExists)
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, domainError.NewCustomError(http.StatusInternalServerError, "failed to hash password", err)
	}

	// Create user entity using shared entity
	user := entity.NewUser(req.Email, req.Username, hashedPassword, req.FirstName, req.LastName)

	// Save user
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, domainError.NewCustomError(http.StatusInternalServerError, "failed to create user", err)
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// Login authenticates a user
func (a *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domainError.NewCustomError(http.StatusUnauthorized, "invalid credentials", domainError.ErrInvalidCredentials)
	}

	// Verify password
	if err := crypto.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, domainError.NewCustomError(http.StatusUnauthorized, "invalid credentials", domainError.ErrInvalidCredentials)
	}

	// Generate token with string UUID
	token, err := jwt.GenerateToken(a.jwtSecret, jwt.UserPayload{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		return nil, domainError.NewCustomError(http.StatusInternalServerError, "failed to generate token", err)
	}

	// Remove password from response
	user.Password = ""

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}
