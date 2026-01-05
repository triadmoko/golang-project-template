package usecase

import (
	"app/internal/features/user/delivery/http/dto"
	mocks "app/internal/mocks/repository"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/domain/entity"
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

func setupTest(t *testing.T) (*userUsecase, *mocks.MockUserRepository) {
	mockRepo := mocks.NewMockUserRepository(t)
	logger := logrus.New()
	logger.SetOutput(os.Stderr)

	uc := &userUsecase{
		userRepo: mockRepo,
		logger:   logger,
	}

	return uc, mockRepo
}

func createTestContext() context.Context {
	return context.WithValue(context.Background(), middleware.LangKey, constants.LangEN)
}

func TestGetProfile_Success(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	userID := "user-123"
	expectedUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(expectedUser, nil)

	user, status, err := uc.GetProfile(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Empty(t, user.Password) // Password should be removed
}

func TestGetProfile_UserNotFound(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	userID := "nonexistent-user"

	mockRepo.EXPECT().GetByID(ctx, userID).Return(nil, errors.New("not found"))

	user, status, err := uc.GetProfile(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, status)
	assert.Nil(t, user)
}

func TestUpdateProfile_Success(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	userID := "user-123"
	existingUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Old",
		LastName:  "Name",
	}

	req := &dto.UpdateProfileRequest{
		FirstName: "New",
		LastName:  "Name",
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	mockRepo.EXPECT().Update(ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, status, err := uc.UpdateProfile(ctx, userID, req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.NotNil(t, user)
	assert.Equal(t, "New", user.FirstName)
	assert.Equal(t, "Name", user.LastName)
	assert.Empty(t, user.Password) // Password should be removed
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	userID := "nonexistent-user"
	req := &dto.UpdateProfileRequest{
		FirstName: "New",
		LastName:  "Name",
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(nil, errors.New("not found"))

	user, status, err := uc.UpdateProfile(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, status)
	assert.Nil(t, user)
}

func TestUpdateProfile_UpdateError(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	userID := "user-123"
	existingUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Old",
		LastName:  "Name",
	}

	req := &dto.UpdateProfileRequest{
		FirstName: "New",
		LastName:  "Name",
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	mockRepo.EXPECT().Update(ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

	user, status, err := uc.UpdateProfile(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Nil(t, user)
}

func TestUpdateProfile_PartialUpdate(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	userID := "user-123"
	existingUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Old",
		LastName:  "Name",
	}

	// Only update FirstName
	req := &dto.UpdateProfileRequest{
		FirstName: "NewFirst",
		LastName:  "",
	}

	mockRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	mockRepo.EXPECT().Update(ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, status, err := uc.UpdateProfile(ctx, userID, req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.NotNil(t, user)
	assert.Equal(t, "NewFirst", user.FirstName)
	assert.Equal(t, "Name", user.LastName) // Should keep original
}

func TestGetUsers_Success(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	limit := 10
	offset := 0

	expectedUsers := []*entity.User{
		{
			ID:        "user-1",
			Email:     "user1@example.com",
			Username:  "user1",
			Password:  "hashedpassword",
			FirstName: "User",
			LastName:  "One",
		},
		{
			ID:        "user-2",
			Email:     "user2@example.com",
			Username:  "user2",
			Password:  "hashedpassword",
			FirstName: "User",
			LastName:  "Two",
		},
	}

	mockRepo.EXPECT().List(ctx, limit, offset).Return(expectedUsers, nil)

	users, status, err := uc.GetUsers(ctx, limit, offset)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Len(t, users, 2)
	// Password should be removed from all users
	for _, user := range users {
		assert.Empty(t, user.Password)
	}
}

func TestGetUsers_Error(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	limit := 10
	offset := 0

	mockRepo.EXPECT().List(ctx, limit, offset).Return(nil, errors.New("database error"))

	users, status, err := uc.GetUsers(ctx, limit, offset)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Nil(t, users)
}

func TestGetUsers_EmptyList(t *testing.T) {
	uc, mockRepo := setupTest(t)
	ctx := createTestContext()

	limit := 10
	offset := 0

	mockRepo.EXPECT().List(ctx, limit, offset).Return([]*entity.User{}, nil)

	users, status, err := uc.GetUsers(ctx, limit, offset)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Empty(t, users)
}
