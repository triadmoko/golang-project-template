package handler

import (
	"app/internal/features/user/delivery/http/dto"
	"app/internal/features/user/usecase"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/delivery/http/response"
	domainError "app/internal/shared/domain/error"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=entity.User}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	user, err := h.userUsecase.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get profile", err)
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} response.SuccessResponse{data=entity.User}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userUsecase.UpdateProfile(c.Request.Context(), userID.(string), &usecase.UpdateProfileRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", user)
}

// GetUsers handles getting list of users
// @Summary Get users list
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.SuccessResponse{data=[]entity.User}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := h.userUsecase.GetUsers(c.Request.Context(), limit, offset)
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}
