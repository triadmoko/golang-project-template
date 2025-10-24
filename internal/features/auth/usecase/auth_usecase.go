package usecase

import (
	"app/internal/features/auth/delivery/http/dto"
	"app/internal/features/auth/domain/entity"
	"app/internal/features/auth/domain/repository"
	"app/internal/features/auth/domain/service"
	domainError "app/internal/shared/domain/error"
	"context"
	"net/http"
)

// AuthUsecase defines the interface for authentication use cases
type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*LoginResponse, error)
}

// authUsecase implements AuthUsecase interface
type authUsecase struct {
	userRepo    repository.UserRepository
	authService service.AuthService
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo repository.UserRepository, authService service.AuthService) AuthUsecase {
	return &authUsecase{
		userRepo:    userRepo,
		authService: authService,
	}
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	User  *entity.User `json:"user"`
	Token string       `json:"token"`
}

// Register creates a new user
func (a *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*entity.User, error) {
	// Check if user already exists
	existingUser, _ := a.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domainError.NewCustomError(http.StatusBadRequest, "user already exists", domainError.ErrUserAlreadyExists)
	}

	existingUser, _ = a.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, domainError.NewCustomError(400, "username already taken", domainError.ErrUserAlreadyExists)
	}

	// Hash password
	hashedPassword, err := a.authService.HashPassword(req.Password)
	if err != nil {
		return nil, domainError.NewCustomError(500, "failed to hash password", err)
	}

	// Create user entity
	user := entity.NewUser(req)
	user.Password = hashedPassword

	// Save user
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, domainError.NewCustomError(500, "failed to create user", err)
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
		return nil, domainError.NewCustomError(401, "invalid credentials", domainError.ErrInvalidCredentials)
	}

	// Verify password
	if err := a.authService.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, domainError.NewCustomError(401, "invalid credentials", domainError.ErrInvalidCredentials)
	}

	// Generate token
	token, err := a.authService.GenerateToken(user)
	if err != nil {
		return nil, domainError.NewCustomError(500, "failed to generate token", err)
	}

	// Remove password from response
	user.Password = ""

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}
