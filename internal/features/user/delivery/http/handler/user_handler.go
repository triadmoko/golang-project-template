package handler

import (
	"app/internal/features/user/delivery/http/dto"
	"app/internal/features/user/usecase"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/delivery/http/response"
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
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	lang := middleware.GetLangFromGin(c)

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.SomethingWentWrong, lang), nil)
		return
	}

	user, status, err := h.userUsecase.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		response.NewResponse(c, status, nil, err.Error(), nil)
		return
	}

	response.NewResponse(c, status, user, "Profile retrieved successfully", nil)
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	lang := middleware.GetLangFromGin(c)

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.SomethingWentWrong, lang), nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponse(c, http.StatusBadRequest, nil, constants.GetErrorMessage(constants.ValidationFailed, lang), map[string][]string{
			"body": {err.Error()},
		})
		return
	}

	// Validate request
	if errors := req.Validate(lang); len(errors) > 0 {
		response.NewResponse(c, http.StatusBadRequest, nil, constants.GetErrorMessage(constants.ValidationFailed, lang), errors)
		return
	}

	user, status, err := h.userUsecase.UpdateProfile(c.Request.Context(), userID.(string), &req)
	if err != nil {
		response.NewResponse(c, status, nil, err.Error(), nil)
		return
	}

	response.NewResponse(c, status, user, "Profile updated successfully", nil)
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
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
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

	users, status, err := h.userUsecase.GetUsers(c.Request.Context(), limit, offset)
	if err != nil {
		response.NewResponse(c, status, nil, err.Error(), nil)
		return
	}

	response.NewResponse(c, status, users, "Users retrieved successfully", nil)
}
