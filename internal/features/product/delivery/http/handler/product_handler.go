package handler

import (
	"app/internal/features/product/delivery/http/dto"
	"app/internal/features/product/usecase"
	"app/internal/shared/delivery/http/response"
	domainError "app/internal/shared/domain/error"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

// CreateProduct handles product creation
// @Summary Create a new product
// @Description Create a new product with name, description, price, stock, and category
// @Tags products
// @Accept json
// @Produce json
// @Param request body usecase.CreateProductRequest true "Product data"
// @Success 201 {object} response.SuccessResponse{data=entity.Product}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	product, err := h.productUsecase.CreateProduct(c.Request.Context(), &usecase.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	})
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

// GetProduct handles getting a product by ID
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.SuccessResponse{data=entity.Product}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	product, err := h.productUsecase.GetProduct(c.Request.Context(), productID)
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get product", err)
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

// UpdateProduct handles updating a product
// @Summary Update product
// @Description Update a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body usecase.UpdateProductRequest true "Product update data"
// @Success 200 {object} response.SuccessResponse{data=entity.Product}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	product, err := h.productUsecase.UpdateProduct(c.Request.Context(), productID, &usecase.UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		IsActive:    req.IsActive,
	})
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update product", err)
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)
}

// DeleteProduct handles deleting a product
// @Summary Delete product
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required", nil)
		return
	}

	err := h.productUsecase.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}

// GetProducts handles getting list of products
// @Summary Get products list
// @Description Get a paginated list of products
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.SuccessResponse{data=[]entity.Product}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
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

	products, err := h.productUsecase.GetProducts(c.Request.Context(), limit, offset)
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get products", err)
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", products)
}

// GetProductsByCategory handles getting products by category
// @Summary Get products by category
// @Description Get a paginated list of products filtered by category
// @Tags products
// @Accept json
// @Produce json
// @Param category path string true "Product category"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.SuccessResponse{data=[]entity.Product}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products/category/{category} [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		response.Error(c, http.StatusBadRequest, "Category is required", nil)
		return
	}

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

	products, err := h.productUsecase.GetProductsByCategory(c.Request.Context(), category, limit, offset)
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get products by category", err)
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", products)
}

// SearchProducts handles searching products
// @Summary Search products
// @Description Search products by name, description, or category
// @Tags products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.SuccessResponse{data=[]entity.Product}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.Error(c, http.StatusBadRequest, "Search query is required", nil)
		return
	}

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

	products, err := h.productUsecase.SearchProducts(c.Request.Context(), query, limit, offset)
	if err != nil {
		if customErr, ok := err.(*domainError.CustomError); ok {
			response.Error(c, customErr.Code, customErr.Message, customErr.Err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to search products", err)
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", products)
}
