package handler

import (
	"app/internal/features/user/delivery/http/dto"
	mocks "app/internal/mocks/usecase"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/domain/entity"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setupGinContext(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func setLanguageMiddleware(c *gin.Context) {
	c.Set(middleware.LangKey, constants.LangEN)
}

func setUserIDMiddleware(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.LangKey, constants.LangEN)
		c.Set(middleware.UserIDKey, userID)
	}
}

func TestGetProfile_Success(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	userID := "user-123"
	router := setupTestRouter()
	router.GET("/profile", setUserIDMiddleware(userID), handler.GetProfile)

	expectedUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	mockUsecase.EXPECT().
		GetProfile(mock.Anything, userID).
		Return(expectedUser, http.StatusOK, nil)

	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["error"].(bool))
}

func TestGetProfile_NoUserID(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	router := setupTestRouter()
	// Not setting user ID in context
	router.GET("/profile", setLanguageMiddleware, handler.GetProfile)

	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProfile_UsecaseError(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	userID := "nonexistent-user"
	router := setupTestRouter()
	router.GET("/profile", setUserIDMiddleware(userID), handler.GetProfile)

	mockUsecase.EXPECT().
		GetProfile(mock.Anything, userID).
		Return(nil, http.StatusNotFound, errors.New("user not found"))

	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateProfile_Success(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	userID := "user-123"
	router := setupTestRouter()
	router.PUT("/profile", setUserIDMiddleware(userID), handler.UpdateProfile)

	reqBody := dto.UpdateProfileRequest{
		FirstName: "NewFirst",
		LastName:  "NewLast",
	}

	expectedUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
	}

	mockUsecase.EXPECT().
		UpdateProfile(mock.Anything, userID, &reqBody).
		Return(expectedUser, http.StatusOK, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["error"].(bool))
}

func TestUpdateProfile_NoUserID(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	router := setupTestRouter()
	router.PUT("/profile", setLanguageMiddleware, handler.UpdateProfile)

	reqBody := dto.UpdateProfileRequest{
		FirstName: "NewFirst",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateProfile_BindJSONError(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	userID := "user-123"
	router := setupTestRouter()
	router.PUT("/profile", setUserIDMiddleware(userID), handler.UpdateProfile)

	req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProfile_ValidationError(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	userID := "user-123"
	router := setupTestRouter()
	router.PUT("/profile", setUserIDMiddleware(userID), handler.UpdateProfile)

	// Empty request - validation should fail
	reqBody := dto.UpdateProfileRequest{
		FirstName: "",
		LastName:  "",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["error"].(bool))
}

func TestUpdateProfile_UsecaseError(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	userID := "user-123"
	router := setupTestRouter()
	router.PUT("/profile", setUserIDMiddleware(userID), handler.UpdateProfile)

	reqBody := dto.UpdateProfileRequest{
		FirstName: "NewFirst",
		LastName:  "NewLast",
	}

	mockUsecase.EXPECT().
		UpdateProfile(mock.Anything, userID, &reqBody).
		Return(nil, http.StatusInternalServerError, errors.New("database error"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetUsers_Success(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	router := setupTestRouter()
	router.GET("/users", setLanguageMiddleware, handler.GetUsers)

	expectedUsers := []*entity.User{
		{
			ID:        "user-1",
			Email:     "user1@example.com",
			Username:  "user1",
			FirstName: "User",
			LastName:  "One",
		},
		{
			ID:        "user-2",
			Email:     "user2@example.com",
			Username:  "user2",
			FirstName: "User",
			LastName:  "Two",
		},
	}

	mockUsecase.EXPECT().
		GetUsers(mock.Anything, 10, 0).
		Return(expectedUsers, http.StatusOK, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["error"].(bool))
}

func TestGetUsers_WithPagination(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	router := setupTestRouter()
	router.GET("/users", setLanguageMiddleware, handler.GetUsers)

	expectedUsers := []*entity.User{
		{
			ID:        "user-3",
			Email:     "user3@example.com",
			Username:  "user3",
			FirstName: "User",
			LastName:  "Three",
		},
	}

	mockUsecase.EXPECT().
		GetUsers(mock.Anything, 5, 10).
		Return(expectedUsers, http.StatusOK, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users?limit=5&offset=10", nil)
	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUsers_InvalidPagination(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	router := setupTestRouter()
	router.GET("/users", setLanguageMiddleware, handler.GetUsers)

	expectedUsers := []*entity.User{}

	// Invalid values should default to 10 and 0
	mockUsecase.EXPECT().
		GetUsers(mock.Anything, 10, 0).
		Return(expectedUsers, http.StatusOK, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users?limit=invalid&offset=-5", nil)
	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUsers_UsecaseError(t *testing.T) {
	mockUsecase := mocks.NewMockUserUsecase(t)
	handler := NewUserHandler(mockUsecase)

	router := setupTestRouter()
	router.GET("/users", setLanguageMiddleware, handler.GetUsers)

	mockUsecase.EXPECT().
		GetUsers(mock.Anything, 10, 0).
		Return(nil, http.StatusInternalServerError, errors.New("database error"))

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
