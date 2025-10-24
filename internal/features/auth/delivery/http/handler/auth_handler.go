package handler

import (
	"app/internal/features/auth/delivery/http/dto"
	"app/internal/features/auth/usecase"
	"app/internal/shared/delivery/http/response"
	domainError "app/internal/shared/domain/error"
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

	user, err := h.authUsecase.Register(c.Request.Context(), &usecase.RegisterRequest{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", user)
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

	loginResp, err := h.authUsecase.Login(c.Request.Context(), &usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", loginResp)
}
