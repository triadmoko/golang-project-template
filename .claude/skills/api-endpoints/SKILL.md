---
name: api-endpoints
description: Create REST API endpoints with proper validation, error handling, and Swagger documentation. Use when adding new routes, HTTP handlers, or REST API functionality.
allowed-tools: Read, Write, Edit, Bash
---

# API Endpoints

This skill guides you through creating REST API endpoints following this project's actual patterns from the codebase.

## ⚠️ CRITICAL PATTERNS

### Use Case Return Pattern
```go
// Use cases return (data, int, error) - int is HTTP status code!
func Create(ctx context.Context, req dto.Request) (*dto.Response, int, error)
```

### Validation Pattern
```go
// DTOs use ONLY json tags, NO binding tags
type Request struct {
    Email string `json:"email"`  // NO binding:"required,email"
}

// Validation via custom method
func (r *Request) Validate(lang constants.Lang) map[string][]string {
    errors := make(map[string][]string)
    if r.Email == "" {
        errors["email"] = append(errors["email"], "email is required")
    }
    return errors
}
```

### Handler Pattern
```go
// 1. Bind JSON
if err := c.ShouldBindJSON(&req); err != nil {
    response.NewResponse(c, http.StatusBadRequest, nil,
        constants.GetErrorMessage(constants.ValidationFailed, lang),
        map[string][]string{"body": {err.Error()}})
    return
}

// 2. Custom validation
if errors := req.Validate(lang); len(errors) > 0 {
    response.NewResponse(c, status, nil, message, errors)
    return
}

// 3. Call use case with context and use status from usecase
result, status, err := h.usecase.Create(c.Request.Context(), req)
if err != nil {
    response.NewResponse(c, status, nil, err.Error(), nil)
    return
}

response.NewResponse(c, status, result, "success", nil)
```

## Endpoint Checklist

- [ ] DTO with json tags only (no binding tags)
- [ ] Custom Validate(lang) method returning map[string][]string
- [ ] Handler binds JSON, then validates
- [ ] Handler passes c.Request.Context() to use case
- [ ] Use case returns (data, int, error) with HTTP status
- [ ] JWT claims type is *jwt.Claims with claims.UserID field
- [ ] Swagger comments complete
- [ ] Route registered in module.go

## Creating an Endpoint

### 1. Define DTOs with Custom Validation

Create request/response structures in `internal/features/{feature}/delivery/http/dto/`:

```go
package dto

import (
    "app/internal/shared/constants"
    "app/internal/shared/domain/entity"
    "fmt"
)

// CreateProductRequest - NO binding tags, only json!
type CreateProductRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    CategoryID  string  `json:"category_id"`
}

// Custom validation method
func (r *CreateProductRequest) Validate(lang constants.Lang) map[string][]string {
    errors := make(map[string][]string)

    // Required field validation
    if r.Name == "" {
        errors["name"] = append(errors["name"],
            fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "name"))
    }

    // Min/Max length validation
    if len(r.Name) < 3 {
        errors["name"] = append(errors["name"], "name must be at least 3 characters")
    }

    if len(r.Name) > 255 {
        errors["name"] = append(errors["name"], "name must be at most 255 characters")
    }

    // Numeric validation
    if r.Price <= 0 {
        errors["price"] = append(errors["price"], "price must be greater than 0")
    }

    // UUID validation (you can use validation helper if available)
    if r.CategoryID == "" {
        errors["category_id"] = append(errors["category_id"],
            fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "category_id"))
    }

    return errors
}

// UpdateProductRequest
type UpdateProductRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

func (r *UpdateProductRequest) Validate(lang constants.Lang) map[string][]string {
    errors := make(map[string][]string)

    if r.Name == "" {
        errors["name"] = append(errors["name"],
            fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "name"))
    }

    if r.Price <= 0 {
        errors["price"] = append(errors["price"], "price must be greater than 0")
    }

    return errors
}

// ProductResponse
type ProductResponse struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    CategoryID  string    `json:"category_id"`
    CreatedAt   time.Time `json:"created_at"`  // time.Time, not string!
}

// Mapper function
func ToProductResponse(product *entity.Product) *ProductResponse {
    return &ProductResponse{
        ID:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        CategoryID:  product.CategoryID,
        CreatedAt:   product.CreatedAt,  // Direct assignment, no .Format()
    }
}

// List response with pagination
type ProductListResponse struct {
    Products   []*ProductResponse    `json:"products"`
    Pagination *pkg.PaginationResponse `json:"pagination"`
}
```

**IMPORTANT:**
- ✅ Only `json` tags (NO `binding` tags!)
- ✅ Custom `Validate(lang constants.Lang) map[string][]string` method
- ✅ Mapper functions to convert entities to responses

### 2. Create Handler Method

Add handler in `internal/features/{feature}/delivery/http/handler/{feature}_handler.go`:

```go
// CreateProduct godoc
// @Summary     Create a new product
// @Description Create a new product with the provided details
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.CreateProductRequest true "Product details"
// @Success     201 {object} response.Response{data=dto.ProductResponse} "Product created"
// @Failure     400 {object} response.Response "Invalid input"
// @Failure     401 {object} response.Response "Unauthorized"
// @Failure     500 {object} response.Response "Internal server error"
// @Router      /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    lang := middleware.GetLangFromGin(c)

    // Get authenticated user from JWT claims
    claimsVal, exists := c.Get("sess")
    if !exists {
        response.NewResponse(c, http.StatusUnauthorized, nil,
            constants.GetErrorMessage(constants.Unauthorized, lang), nil)
        return
    }

    claims, ok := claimsVal.(*jwt.Claims)  // Type: *jwt.Claims (not *jwt.JWTClaims!)
    if !ok {
        response.NewResponse(c, http.StatusUnauthorized, nil,
            constants.GetErrorMessage(constants.Unauthorized, lang), nil)
        return
    }

    // 1. Bind JSON first
    var req dto.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang),
            map[string][]string{"body": {err.Error()}})  // Structured error!
        return
    }

    // 2. Then custom validation
    if errors := req.Validate(lang); len(errors) > 0 {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang), errors)
        return
    }

    // 3. Call use case with context - returns (data, status, error)
    product, status, err := h.productUsecase.Create(c.Request.Context(), claims.UserID, req)
    if err != nil {
        response.NewResponse(c, status, nil, err.Error(), nil)
        return
    }

    response.NewResponse(c, status, product, "Product created successfully", nil)
}
```

**CRITICAL HANDLER PATTERNS:**
- ✅ Get lang from Gin: `middleware.GetLangFromGin(c)`
- ✅ Claims type: `*jwt.Claims` with field `claims.UserID` (not `claims.ID`)
- ✅ Bind JSON first: `c.ShouldBindJSON(&req)`
- ✅ Then validate: `req.Validate(lang)`
- ✅ Errors as map: `map[string][]string{"body": {err.Error()}}`
- ✅ Pass context: `h.usecase.Create(c.Request.Context(), ...)`
- ✅ Get status from usecase: `product, status, err := ...`

## Swagger Annotations

### Handler Comments Format

```go
// MethodName godoc
// @Summary     Short description
// @Description Detailed description
// @Tags        tag-name
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       name location type required "description"
// @Success     code {object} response.Response{data=dto.ResponseType} "description"
// @Failure     code {object} response.Response "description"
// @Router      /api/v1/path [method]
```

### Common Annotations

**Parameters**:
```go
// Path parameter
// @Param id path string true "Product ID"

// Query parameter
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)

// Body parameter
// @Param request body dto.CreateProductRequest true "Request body"
```

**Success Response with Data**:
```go
// @Success 200 {object} response.Response{data=dto.ProductResponse} "Product retrieved"
// @Success 200 {object} response.Response{data=dto.ProductListResponse} "Product list"
```

## Common Endpoint Patterns

### GET - Single Item

```go
// GetProduct godoc
// @Summary     Get product by ID
// @Description Retrieve a single product by its ID
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Product ID"
// @Success     200 {object} response.Response{data=dto.ProductResponse}
// @Failure     404 {object} response.Response "Product not found"
// @Router      /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
    id := c.Param("id")

    // Call use case with context - returns (data, status, error)
    product, status, err := h.productUsecase.GetByID(c.Request.Context(), id)
    if err != nil {
        response.NewResponse(c, status, nil, err.Error(), nil)
        return
    }

    response.NewResponse(c, status, product, "Product retrieved successfully", nil)
}
```

### GET - List with Pagination

```go
// GetUsers godoc
// @Summary     Get users list
// @Description Get a paginated and filtered list of users
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       per_page query int false "Items per page" default(10)
// @Param       page query int false "Page number" default(1)
// @Param       email query string false "Filter by email"
// @Success     200 {object} response.Response{data=dto.UserListResponse}
// @Router      /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
    // Build queries map from query parameters
    queries := map[string]string{}

    err := c.BindQuery(&queries)
    if err != nil {
        lang := middleware.GetLangFromGin(c)
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang),
            map[string][]string{"query": {err.Error()}})
        return
    }

    // Call use case - returns (data, pagination, status, error)
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
```

### POST - Create

```go
// CreateProduct godoc
// @Summary     Create product
// @Description Create a new product
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.CreateProductRequest true "Product details"
// @Success     201 {object} response.Response{data=dto.ProductResponse}
// @Failure     400 {object} response.Response
// @Router      /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    lang := middleware.GetLangFromGin(c)

    claimsVal, exists := c.Get("sess")
    if !exists {
        response.NewResponse(c, http.StatusUnauthorized, nil,
            constants.GetErrorMessage(constants.Unauthorized, lang), nil)
        return
    }

    claims, ok := claimsVal.(*jwt.Claims)
    if !ok {
        response.NewResponse(c, http.StatusUnauthorized, nil,
            constants.GetErrorMessage(constants.Unauthorized, lang), nil)
        return
    }

    var req dto.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang),
            map[string][]string{"body": {err.Error()}})
        return
    }

    if errors := req.Validate(lang); len(errors) > 0 {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang), errors)
        return
    }

    product, status, err := h.productUsecase.Create(c.Request.Context(), claims.UserID, req)
    if err != nil {
        response.NewResponse(c, status, nil, err.Error(), nil)
        return
    }

    response.NewResponse(c, status, product, "Product created successfully", nil)
}
```

### PUT - Update

```go
// UpdateProduct godoc
// @Summary     Update product
// @Description Update an existing product
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Product ID"
// @Param       request body dto.UpdateProductRequest true "Updated product details"
// @Success     200 {object} response.Response{data=dto.ProductResponse}
// @Failure     400 {object} response.Response
// @Failure     404 {object} response.Response
// @Router      /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
    lang := middleware.GetLangFromGin(c)
    id := c.Param("id")

    var req dto.UpdateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang),
            map[string][]string{"body": {err.Error()}})
        return
    }

    if errors := req.Validate(lang); len(errors) > 0 {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang), errors)
        return
    }

    product, status, err := h.productUsecase.Update(c.Request.Context(), id, req)
    if err != nil {
        response.NewResponse(c, status, nil, err.Error(), nil)
        return
    }

    response.NewResponse(c, status, product, "Product updated successfully", nil)
}
```

### DELETE

```go
// DeleteProduct godoc
// @Summary     Delete product
// @Description Delete a product by ID
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Product ID"
// @Success     200 {object} response.Response
// @Failure     404 {object} response.Response
// @Router      /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
    id := c.Param("id")

    status, err := h.productUsecase.Delete(c.Request.Context(), id)
    if err != nil {
        response.NewResponse(c, status, nil, err.Error(), nil)
        return
    }

    response.NewResponse(c, status, nil, "Product deleted successfully", nil)
}
```

## Custom Validation Examples

### Email Validation
```go
func (r *Request) Validate(lang constants.Lang) map[string][]string {
    errors := make(map[string][]string)

    if r.Email == "" {
        errors["email"] = append(errors["email"],
            fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "email"))
    } else if !constants.IsValidEmail(r.Email) {
        errors["email"] = append(errors["email"],
            constants.GetValidationMessage(constants.InvalidEmail, lang))
    }

    return errors
}
```

### String Length Validation
```go
if r.Username == "" {
    errors["username"] = append(errors["username"],
        fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "username"))
} else {
    if !constants.MinLength(r.Username, 3) {
        errors["username"] = append(errors["username"],
            fmt.Sprintf(constants.GetValidationMessage(constants.UsernameTooShort, lang), 3))
    }
    if !constants.MaxLength(r.Username, 20) {
        errors["username"] = append(errors["username"],
            fmt.Sprintf(constants.GetValidationMessage(constants.UsernameTooLong, lang), 20))
    }
}
```

### Numeric Validation
```go
if r.Price <= 0 {
    errors["price"] = append(errors["price"], "price must be greater than 0")
}

if r.Quantity < 1 || r.Quantity > 1000 {
    errors["quantity"] = append(errors["quantity"], "quantity must be between 1 and 1000")
}
```

## Register Routes in Module

Add routes in `module.go`:

```go
func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
    // Public routes (no auth required)
    productsPublic := rg.Group("/products")
    {
        productsPublic.GET("", m.handler.ListProducts)
        productsPublic.GET("/:id", m.handler.GetProduct)
    }

    // Protected routes (auth required)
    productsProtected := rg.Group("/products")
    productsProtected.Use(middleware.AuthMiddleware())
    {
        productsProtected.POST("", m.handler.CreateProduct)
        productsProtected.PUT("/:id", m.handler.UpdateProduct)
        productsProtected.DELETE("/:id", m.handler.DeleteProduct)
    }
}
```

## Regenerate Swagger Docs

After adding or modifying endpoints:

```bash
make swag
```

This updates `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`.

Access documentation at: `http://localhost:8080/swagger/index.html`

## Testing Endpoints

### Using curl

```bash
# Login first
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.data.token')

# Use token in requests
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN"

# Create resource
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Product","price":99.99}'
```

## Error Handling Best Practices

1. **Use language-specific errors from use case**:
   ```go
   // Use case returns error with language support
   product, status, err := h.usecase.Create(c.Request.Context(), req)
   if err != nil {
       response.NewResponse(c, status, nil, err.Error(), nil)
       return
   }
   ```

2. **HTTP status codes from use case**:
   - Use case returns status code
   - Handler uses it directly in response

3. **Structured validation errors**:
   ```go
   // Bind errors
   map[string][]string{"body": {err.Error()}}

   // Custom validation errors
   map[string][]string{
       "email": {"email is required", "invalid email format"},
       "price": {"price must be greater than 0"},
   }
   ```

4. **Unified response format**:
   ```go
   response.NewResponse(c, statusCode, data, message, errors)
   ```

## Quick Reference

**Critical Patterns:**
- ✅ DTOs with json tags ONLY (no binding tags)
- ✅ Custom Validate(lang) method returning map[string][]string
- ✅ Handler: Bind → Validate → Call use case with context
- ✅ Use case returns (data, int, error) with HTTP status
- ✅ JWT claims: `*jwt.Claims` type with `claims.UserID` field
- ✅ Pass context: `h.usecase.Method(c.Request.Context(), ...)`
- ✅ Structured errors: `map[string][]string`
