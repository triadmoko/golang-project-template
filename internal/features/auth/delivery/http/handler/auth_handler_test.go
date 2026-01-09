package handler

import (
	authdto "app/internal/features/auth/delivery/http/dto"
	mocks "app/internal/mocks/usecase"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
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

func TestRegister_Success(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/register", setLanguageMiddleware, handler.Register)

	reqBody := authdto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	expectedUser := &authdto.RegisterResponse{
		ID:        "user-123",
		Email:     reqBody.Email,
		Username:  reqBody.Username,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
	}

	mockUsecase.EXPECT().
		Register(mock.Anything, reqBody).
		Return(expectedUser, http.StatusCreated, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["error"].(bool))
}

func TestRegister_BindJSONError(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/register", setLanguageMiddleware, handler.Register)

	// Invalid JSON body
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_ValidationError(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/register", setLanguageMiddleware, handler.Register)

	// Missing required fields
	reqBody := authdto.RegisterRequest{
		Email: "", // Empty email
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["error"].(bool))
	assert.NotNil(t, response["errors"])
}

func TestRegister_UsecaseError(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/register", setLanguageMiddleware, handler.Register)

	reqBody := authdto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	mockUsecase.EXPECT().
		Register(mock.Anything, reqBody).
		Return(nil, http.StatusBadRequest, errors.New("user already exists"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["error"].(bool))
}

func TestLogin_Success(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/login", setLanguageMiddleware, handler.Login)

	reqBody := authdto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResponse := &authdto.LoginResponse{
		User: &authdto.RegisterResponse{
			ID:       "user-123",
			Email:    reqBody.Email,
			Username: "testuser",
		},
		Token: "jwt-token-here",
	}

	mockUsecase.EXPECT().
		Login(mock.Anything, reqBody).
		Return(expectedResponse, http.StatusOK, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["error"].(bool))
}

func TestLogin_BindJSONError(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/login", setLanguageMiddleware, handler.Login)

	// Invalid JSON body
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_ValidationError(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/login", setLanguageMiddleware, handler.Login)

	// Missing required fields
	reqBody := authdto.LoginRequest{
		Email: "", // Empty email
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["error"].(bool))
	assert.NotNil(t, response["errors"])
}

func TestLogin_UsecaseError(t *testing.T) {
	mockUsecase := mocks.NewMockAuthUsecase(t)
	handler := NewAuthHandler(mockUsecase)

	router := setupTestRouter()
	router.POST("/login", setLanguageMiddleware, handler.Login)

	reqBody := authdto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockUsecase.EXPECT().
		Login(mock.Anything, reqBody).
		Return(nil, http.StatusUnauthorized, errors.New("invalid credentials"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := setupGinContext(router, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["error"].(bool))
}
