package handler

import (
	"app/internal/features/auth/delivery/http/dto"
	"app/internal/features/auth/usecase"
	"app/internal/shared/delivery/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, username, password, first name, and last name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} response.SuccessResponse{data=entity.User}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, status, err := h.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, status, err.Error(), nil)
		return
	}

	response.Success(c, status, "User registered successfully", user)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login data"
// @Success 200 {object} response.SuccessResponse{data=dto.LoginResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	loginResp, status, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, status, err.Error(), nil)
		return
	}

	response.Success(c, status, "Login successful", loginResp)
}
