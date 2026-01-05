package usecase

import (
	"app/internal/features/auth/delivery/http/dto"
	mocks "app/internal/mocks/repository"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/domain/entity"
	"app/pkg/crypto"
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Set JWT_SECRET for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	code := m.Run()
	os.Exit(code)
}

func setupTest(t *testing.T) (*authUsecase, *mocks.MockUserRepository) {
	mockRepo := mocks.NewMockUserRepository(t)
	logger := logrus.New()
	logger.SetOutput(os.Stderr)

	uc := &authUsecase{
		userRepo: mockRepo,
		logger:   logger,
	}

	return uc, mockRepo
}

func createTestContext() context.Context {
	return context.WithValue(context.Background(), middleware.LangKey, constants.LangEN)
}

func TestRegister_Success(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	req := dto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock: email not found
	mockRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("not found"))
	// Mock: username not found
	mockRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, errors.New("not found"))
	// Mock: create user success
	mockRepo.EXPECT().Create(ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, status, err := uc.Register(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Username, user.Username)
	assert.Empty(t, user.Password) // Password should be removed from response
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	req := dto.RegisterRequest{
		Email:     "existing@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	existingUser := &entity.User{
		ID:    "existing-id",
		Email: req.Email,
	}

	// Mock: email already exists
	mockRepo.EXPECT().GetByEmail(ctx, req.Email).Return(existingUser, nil)

	user, status, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, status)
	assert.Nil(t, user)
}

func TestRegister_UsernameAlreadyTaken(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	req := dto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "existinguser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	existingUser := &entity.User{
		ID:       "existing-id",
		Username: req.Username,
	}

	// Mock: email not found
	mockRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("not found"))
	// Mock: username already exists
	mockRepo.EXPECT().GetByUsername(ctx, req.Username).Return(existingUser, nil)

	user, status, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, status)
	assert.Nil(t, user)
}

func TestRegister_CreateUserError(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	req := dto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock: email not found
	mockRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("not found"))
	// Mock: username not found
	mockRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, errors.New("not found"))
	// Mock: create user fails
	mockRepo.EXPECT().Create(ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

	user, status, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Nil(t, user)
}

func TestLogin_Success(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	password := "password123"
	hashedPassword, err := crypto.HashPassword(password)
	require.NoError(t, err)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	existingUser := &entity.User{
		ID:       "user-123",
		Email:    req.Email,
		Username: "testuser",
		Password: hashedPassword,
	}

	// Mock: get user by email success
	mockRepo.EXPECT().GetByEmail(ctx, req.Email).Return(existingUser, nil)

	loginResp, status, err := uc.Login(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.NotNil(t, loginResp)
	assert.NotEmpty(t, loginResp.Token)
	assert.Empty(t, loginResp.User.Password) // Password should be removed
}

func TestLogin_UserNotFound(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	req := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Mock: user not found
	mockRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("not found"))

	loginResp, status, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Nil(t, loginResp)
}

func TestLogin_InvalidPassword(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	correctPassword := "correctpassword"
	hashedPassword, err := crypto.HashPassword(correctPassword)
	require.NoError(t, err)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	existingUser := &entity.User{
		ID:       "user-123",
		Email:    req.Email,
		Username: "testuser",
		Password: hashedPassword,
	}

	// Mock: get user by email success
	mockRepo.EXPECT().GetByEmail(ctx, req.Email).Return(existingUser, nil)

	loginResp, status, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Nil(t, loginResp)
}
