---
name: adding-features
description: Guide for adding new feature modules following Clean Architecture pattern. Use when creating new features, modules, REST endpoints, or business logic components in this Go project.
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

# Adding Features

This skill guides you through adding new feature modules following this project's Clean Architecture and feature-based modular design with **correct patterns** from the actual codebase.

## Feature Module Structure

Each feature follows this structure:

```
internal/features/myfeature/
├── module.go                    # DI container + route registration
├── delivery/http/
│   ├── handler/
│   │   └── myfeature_handler.go # HTTP handlers
│   └── dto/
│       └── myfeature_dto.go     # Request/response DTOs + validation
└── usecase/
    └── myfeature_usecase.go     # Business logic + interface
```

## Critical Patterns from Actual Code

### ⚠️ IMPORTANT: Context Pattern
- **ALL** use case methods MUST accept `ctx context.Context` as first parameter
- **ALL** repository methods MUST accept `ctx context.Context` as first parameter
- Handlers pass `c.Request.Context()` to use cases

### ⚠️ IMPORTANT: Use Case Return Pattern
Use cases return `(data, int, error)` where int is HTTP status code:
```go
func (u *myUsecase) Create(ctx context.Context, req dto.Request) (*dto.Response, int, error)
```

### ⚠️ IMPORTANT: Validation Pattern
- DTOs do NOT use binding tags (only json tags)
- Validation is done via custom `Validate(lang constants.Lang) map[string][]string` method
- Handler calls `ShouldBindJSON()` first, then `req.Validate(lang)`

### ⚠️ IMPORTANT: ID Type
- User.ID is `string` (not `uuid.UUID`)
- UUID stored as string in database

## Step-by-Step Process

### 1. Create Feature Directory Structure

```bash
mkdir -p internal/features/{feature-name}/delivery/http/{handler,dto}
mkdir -p internal/features/{feature-name}/usecase
```

### 2. Create Domain Entities (if needed)

If your feature needs new entities, add them to `internal/shared/domain/entity/`:

```go
// internal/shared/domain/entity/product.go
package entity

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Product struct {
    ID          string         `json:"id" gorm:"type:varchar(36);primaryKey"`
    Name        string         `json:"name" gorm:"type:varchar(255);not null"`
    Description string         `json:"description" gorm:"type:text"`
    Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
    UserID      string         `json:"user_id" gorm:"type:varchar(36);not null"`
    CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Product) TableName() string {
    return "products"
}

// Constructor with UUID generation
func NewProduct(name, description string, price float64, userID string) *Product {
    return &Product{
        ID:          uuid.New().String(),
        Name:        name,
        Description: description,
        Price:       price,
        UserID:      userID,
    }
}

// BeforeCreate hook
func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.ID == "" {
        p.ID = uuid.New().String()
    }
    return nil
}
```

### 3. Create Repository Interface (if needed)

Add repository interfaces to `internal/shared/domain/repository/`:

```go
// internal/shared/domain/repository/product_repository.go
package repository

import (
    "app/internal/shared/domain/entity"
    "context"
)

type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    GetByID(ctx context.Context, id string) (*entity.Product, error)
    Update(ctx context.Context, product *entity.Product) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, userID string, limit, offset int) ([]*entity.Product, int, error)
}
```

**IMPORTANT:** All methods accept `ctx context.Context` as first parameter!

### 4. Implement Repository

Create implementation in `internal/shared/infrastructure/repository/`:

```go
// internal/shared/infrastructure/repository/product_repository_impl.go
package repository

import (
    "app/internal/shared/domain/entity"
    "app/internal/shared/domain/repository"
    "context"
    "gorm.io/gorm"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
    return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
    var product entity.Product
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
        return nil, err
    }
    return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
    return r.db.WithContext(ctx).Updates(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Product{}).Error
}

func (r *productRepository) List(ctx context.Context, userID string, limit, offset int) ([]*entity.Product, int, error) {
    var products []*entity.Product
    query := r.db.WithContext(ctx).Where("user_id = ?", userID)

    var total int64
    if err := query.Model(&entity.Product{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
        return nil, 0, err
    }

    return products, int(total), nil
}
```

### 5. Create DTOs with Validation

Define request/response structures in `internal/features/{feature-name}/delivery/http/dto/`:

```go
// internal/features/product/delivery/http/dto/product_dto.go
package dto

import (
    "app/internal/shared/constants"
    "app/internal/shared/domain/entity"
    "fmt"
)

// CreateProductRequest - NO binding tags!
type CreateProductRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

// Validate method for custom validation
func (r *CreateProductRequest) Validate(lang constants.Lang) map[string][]string {
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
    UserID      string    `json:"user_id"`
    CreatedAt   time.Time `json:"created_at"`  // time.Time, not string!
}

// Mapper function
func ToProductResponse(product *entity.Product) *ProductResponse {
    return &ProductResponse{
        ID:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        UserID:      product.UserID,
        CreatedAt:   product.CreatedAt,  // Direct assignment, no .Format()
    }
}
```

**IMPORTANT:**
- DTOs use ONLY `json` tags, NO `binding` tags
- Validation via custom `Validate(lang)` method
- Returns `map[string][]string` for field-specific errors
- Include mapper functions like `ToProductResponse`

### 6. Create Use Case

Define interface and implementation in `internal/features/{feature-name}/usecase/`:

```go
// internal/features/product/usecase/product_usecase.go
package usecase

import (
    "app/internal/features/product/delivery/http/dto"
    "app/internal/shared/constants"
    "app/internal/shared/delivery/http/middleware"
    "app/internal/shared/domain/entity"
    "app/internal/shared/domain/repository"
    "context"
    "net/http"

    "github.com/sirupsen/logrus"
)

// Interface
type ProductUsecase interface {
    Create(ctx context.Context, userID string, req dto.CreateProductRequest) (*dto.ProductResponse, int, error)
    GetByID(ctx context.Context, id string) (*dto.ProductResponse, int, error)
    Update(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, int, error)
    Delete(ctx context.Context, id string) (int, error)
}

// Implementation
type productUsecase struct {
    productRepo repository.ProductRepository
    logger      *logrus.Logger
}

func NewProductUsecase(productRepo repository.ProductRepository, logger *logrus.Logger) ProductUsecase {
    return &productUsecase{
        productRepo: productRepo,
        logger:      logger,
    }
}

// Create product
func (u *productUsecase) Create(ctx context.Context, userID string, req dto.CreateProductRequest) (*dto.ProductResponse, int, error) {
    // Get language from context
    lang := middleware.GetLangFromContext(ctx)

    // Create entity using constructor
    product := entity.NewProduct(req.Name, req.Description, req.Price, userID)

    // Save to repository
    if err := u.productRepo.Create(ctx, product); err != nil {
        u.logger.Error("u.productRepo.Create ", err)
        return nil, http.StatusInternalServerError, constants.GetError(constants.SomethingWentWrong, lang)
    }

    // Return with mapper
    return dto.ToProductResponse(product), http.StatusCreated, nil
}

// GetByID retrieves product
func (u *productUsecase) GetByID(ctx context.Context, id string) (*dto.ProductResponse, int, error) {
    lang := middleware.GetLangFromContext(ctx)

    product, err := u.productRepo.GetByID(ctx, id)
    if err != nil {
        u.logger.Error("u.productRepo.GetByID ", err)
        return nil, http.StatusNotFound, constants.GetError(constants.UserNotFound, lang) // Use appropriate error constant
    }

    return dto.ToProductResponse(product), http.StatusOK, nil
}

// Update product
func (u *productUsecase) Update(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, int, error) {
    lang := middleware.GetLangFromContext(ctx)

    // Get existing product
    product, err := u.productRepo.GetByID(ctx, id)
    if err != nil {
        u.logger.Error("u.productRepo.GetByID ", err)
        return nil, http.StatusNotFound, constants.GetError(constants.UserNotFound, lang)
    }

    // Update fields
    product.Name = req.Name
    product.Description = req.Description
    product.Price = req.Price

    // Save
    if err := u.productRepo.Update(ctx, product); err != nil {
        u.logger.Error("u.productRepo.Update ", err)
        return nil, http.StatusInternalServerError, constants.GetError(constants.SomethingWentWrong, lang)
    }

    return dto.ToProductResponse(product), http.StatusOK, nil
}

// Delete product
func (u *productUsecase) Delete(ctx context.Context, id string) (int, error) {
    lang := middleware.GetLangFromContext(ctx)

    // Check if exists
    _, err := u.productRepo.GetByID(ctx, id)
    if err != nil {
        u.logger.Error("u.productRepo.GetByID ", err)
        return http.StatusNotFound, constants.GetError(constants.UserNotFound, lang)
    }

    // Delete
    if err := u.productRepo.Delete(ctx, id); err != nil {
        u.logger.Error("u.productRepo.Delete ", err)
        return http.StatusInternalServerError, constants.GetError(constants.SomethingWentWrong, lang)
    }

    return http.StatusOK, nil
}
```

**CRITICAL PATTERNS:**
- ✅ Methods accept `ctx context.Context` as first parameter
- ✅ Return `(data, int, error)` with HTTP status code
- ✅ Get language via `middleware.GetLangFromContext(ctx)`
- ✅ Use entity constructors: `entity.NewProduct(...)`
- ✅ Use mapper functions: `dto.ToProductResponse(product)`
- ✅ Pass context to repository: `u.productRepo.Create(ctx, product)`

### 7. Create Handler

Implement HTTP handlers in `internal/features/{feature-name}/delivery/http/handler/`:

```go
// internal/features/product/delivery/http/handler/product_handler.go
package handler

import (
    "app/internal/features/product/delivery/http/dto"
    "app/internal/features/product/usecase"
    "app/internal/shared/constants"
    "app/internal/shared/delivery/http/middleware"
    "app/internal/shared/delivery/http/response"
    "app/pkg/jwt"
    "net/http"

    "github.com/gin-gonic/gin"
)

type ProductHandler struct {
    productUsecase usecase.ProductUsecase
}

func NewProductHandler(uc usecase.ProductUsecase) *ProductHandler {
    return &ProductHandler{productUsecase: uc}
}

// CreateProduct godoc
// @Summary     Create product
// @Description Create a new product
// @Tags        products
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body dto.CreateProductRequest true "Product data"
// @Success     201 {object} response.Response{data=dto.ProductResponse}
// @Failure     400 {object} response.Response
// @Failure     401 {object} response.Response
// @Failure     500 {object} response.Response
// @Router      /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    lang := middleware.GetLangFromGin(c)

    // Get user from JWT claims
    claimsVal, exists := c.Get("sess")
    if !exists {
        response.NewResponse(c, http.StatusUnauthorized, nil,
            constants.GetErrorMessage(constants.Unauthorized, lang), nil)
        return
    }

    claims, ok := claimsVal.(*jwt.Claims)  // Type: *jwt.Claims
    if !ok {
        response.NewResponse(c, http.StatusUnauthorized, nil,
            constants.GetErrorMessage(constants.Unauthorized, lang), nil)
        return
    }

    // Bind JSON
    var req dto.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang),
            map[string][]string{"body": {err.Error()}})
        return
    }

    // Custom validation
    if errors := req.Validate(lang); len(errors) > 0 {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.ValidationFailed, lang), errors)
        return
    }

    // Call use case with context
    product, status, err := h.productUsecase.Create(c.Request.Context(), claims.UserID, req)
    if err != nil {
        response.NewResponse(c, status, nil, err.Error(), nil)
        return
    }

    response.NewResponse(c, status, product, "Product created successfully", nil)
}

// GetProduct godoc
// @Summary     Get product
// @Description Get product by ID
// @Tags        products
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Product ID"
// @Success     200 {object} response.Response{data=dto.ProductResponse}
// @Failure     404 {object} response.Response
// @Router      /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
    id := c.Param("id")

    product, status, err := h.productUsecase.GetByID(c.Request.Context(), id)
    if err != nil {
        response.NewResponse(c, status, nil, err.Error(), nil)
        return
    }

    response.NewResponse(c, status, product, "Product retrieved successfully", nil)
}
```

**CRITICAL HANDLER PATTERNS:**
- ✅ Get lang: `middleware.GetLangFromGin(c)`
- ✅ Claims type: `*jwt.Claims` with field `claims.UserID`
- ✅ Bind first: `c.ShouldBindJSON(&req)`
- ✅ Then validate: `req.Validate(lang)`
- ✅ Structured errors: `map[string][]string{"body": {err.Error()}}`
- ✅ Pass context: `h.usecase.Create(c.Request.Context(), ...)`
- ✅ Use status from usecase: `product, status, err := h.usecase.Create(...)`

### 8. Create Module

Wire dependencies in `internal/features/{feature-name}/module.go`:

```go
// internal/features/product/module.go
package product

import (
    "app/internal/features/product/delivery/http/handler"
    "app/internal/features/product/usecase"
    "app/internal/shared/delivery/http/middleware"
    "app/internal/shared/domain/repository"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type Module struct {
    handler *handler.ProductHandler
}

func NewModule(productRepo repository.ProductRepository, logger *logrus.Logger) *Module {
    // Wire dependencies
    uc := usecase.NewProductUsecase(productRepo, logger)
    h := handler.NewProductHandler(uc)

    return &Module{handler: h}
}

func (m *Module) Name() string {
    return "product"
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
    // Protected routes
    productGroup := rg.Group("/products")
    productGroup.Use(middleware.AuthMiddleware())
    {
        productGroup.POST("", m.handler.CreateProduct)
        productGroup.GET("/:id", m.handler.GetProduct)
        productGroup.PUT("/:id", m.handler.UpdateProduct)
        productGroup.DELETE("/:id", m.handler.DeleteProduct)
    }
}
```

### 9. Register Module in App

Add the module to `internal/app/app.go`:

```go
// Initialize repository
productRepo := sharedRepo.NewProductRepository(a.DB.GetDB())

// Register all features
features := []Feature{
    auth.NewModule(userRepo, a.Logger),
    user.NewModule(userRepo, a.Logger),
    product.NewModule(productRepo, a.Logger), // Add your new feature
}
```

### 10. Update Swagger Documentation

```bash
make swag
```

## Quick Reference Checklist

When adding a new feature:

- [ ] Context parameter in ALL use case and repository methods
- [ ] Use case returns `(data, int, error)` with HTTP status code
- [ ] DTOs without binding tags, only json tags
- [ ] Custom `Validate(lang)` method in DTOs returning `map[string][]string`
- [ ] Entity ID is `string` type, not `uuid.UUID`
- [ ] Entity has constructor function `NewEntity(...)`
- [ ] Mapper functions in DTO package: `ToEntityResponse(entity)`
- [ ] Handler binds JSON first, then calls custom validate
- [ ] Handler passes `c.Request.Context()` to use case
- [ ] Handler uses `*jwt.Claims` type with `claims.UserID` field
- [ ] Use case gets lang via `middleware.GetLangFromContext(ctx)`
- [ ] Error response supports both string and `map[string][]string`
- [ ] Generate Swagger docs after changes
