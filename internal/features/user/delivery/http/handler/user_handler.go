package handler

import (
	"app/internal/features/user/delivery/http/dto"
	"app/internal/features/user/usecase"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/delivery/http/response"
	"app/pkg/jwt"
	"net/http"

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
//
//	@Summary		Get user profile
//	@Description	Get the authenticated user's profile information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=dto.UserResponse}
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	lang := middleware.GetLangFromGin(c)

	// Get claims from context
	claimsVal, exists := c.Get("sess")
	if !exists {
		response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
		return
	}

	claims, ok := claimsVal.(*jwt.Claims)
	if !ok {
		response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
		return
	}

	user, status, err := h.userUsecase.GetProfile(c.Request.Context(), claims.UserID)
	if err != nil {
		response.NewResponse(c, status, nil, err.Error(), nil)
		return
	}

	response.NewResponse(c, status, user, "Profile retrieved successfully", nil)
}

// UpdateProfile handles updating user profile
//
//	@Summary		Update user profile
//	@Description	Update the authenticated user's profile information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpdateProfileRequest	true	"Profile update data"
//	@Success		200		{object}	response.Response{data=dto.UserResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	lang := middleware.GetLangFromGin(c)

	// Get claims from context
	claimsVal, exists := c.Get("sess")
	if !exists {
		response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
		return
	}

	claims, ok := claimsVal.(*jwt.Claims)
	if !ok {
		response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
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

	user, status, err := h.userUsecase.UpdateProfile(c.Request.Context(), claims.UserID, &req)
	if err != nil {
		response.NewResponse(c, status, nil, err.Error(), nil)
		return
	}

	response.NewResponse(c, status, user, "Profile updated successfully", nil)
}

// GetUsers handles getting list of users
//
//	@Summary		Get users list
//	@Description	Get a paginated and filtered list of users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			per_page	query		int		false	"Items per page"	default(10)
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			id			query		string	false	"Filter by user ID"
//	@Param			email		query		string	false	"Filter by email"
//	@Param			username	query		string	false	"Filter by username"
//	@Param			first_name	query		string	false	"Filter by first name (LIKE search)"
//	@Param			last_name	query		string	false	"Filter by last name (LIKE search)"
//	@Param			status		query		string	false	"Filter by status"
//	@Param			gender		query		string	false	"Filter by gender"
//	@Param			role		query		string	false	"Filter by role"
//	@Param			provider	query		string	false	"Filter by provider"
//	@Param			genders		query		string	false	"Filter by multiple genders (comma-separated)"
//	@Param			roles		query		string	false	"Filter by multiple roles (comma-separated)"
//	@Success		200			{object}	response.Response{data=dto.UserListResponse}
//	@Failure		400			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		500			{object}	response.Response
//	@Router			/api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Build queries map from query parameters
	queries := map[string]string{}

	err := c.BindQuery(&queries)
	if err != nil {
		lang := middleware.GetLangFromGin(c)
		response.NewResponse(c, http.StatusBadRequest, nil, constants.GetErrorMessage(constants.ValidationFailed, lang), map[string][]string{
			"query": {err.Error()},
		})
		return
	}

	users, pagination, status, err := h.userUsecase.GetUsers(c.Request.Context(), queries)
	if err != nil {
		response.NewResponse(c, status, nil, err.Error(), nil)
		return
	}

	// Build response with DTO
	responseData := dto.UserListResponse{
		Users:      users,
		Pagination: pagination,
	}

	response.NewResponse(c, status, responseData, "Users retrieved successfully", nil)
}
