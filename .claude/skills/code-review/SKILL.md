---
name: code-review
description: Review Go code for this project following Clean Architecture principles, Go best practices, and project conventions. Use when reviewing code, pull requests, or checking code quality.
allowed-tools: Read, Glob, Grep, Bash
---

# Code Review

This skill guides you through reviewing code for this Go project, focusing on Clean Architecture, Go best practices, and project-specific conventions.

## Review Checklist

### Architecture & Structure

- [ ] Feature follows the module pattern (module.go, handler, usecase, dto)
- [ ] Dependencies flow inward (handler → usecase → repository interface)
- [ ] Use cases depend on interfaces, not implementations
- [ ] Domain entities are in `internal/shared/domain/entity/`
- [ ] Repository interfaces in `internal/shared/domain/repository/`
- [ ] Repository implementations in `internal/shared/infrastructure/repository/`

### API & Handlers

- [ ] Handlers use unified response format: `response.NewResponse()`
- [ ] Language-specific errors used: `constants.GetErrorMessage(code, lang)`
- [ ] Request validation uses binding tags
- [ ] Swagger comments complete and accurate
- [ ] Protected routes use `middleware.AuthMiddleware()`
- [ ] User claims retrieved properly from context
- [ ] Appropriate HTTP status codes used

### Business Logic

- [ ] Business logic in use cases, not handlers
- [ ] Use cases are testable (depend on interfaces)
- [ ] Error handling is comprehensive
- [ ] Logging added for important operations
- [ ] No database logic in use cases (only repository calls)

### Data & DTOs

- [ ] DTOs have `// @name` comments for Swagger
- [ ] Validation tags on all required fields
- [ ] Password fields use `json:"-"` tag
- [ ] UUIDs used for IDs (not integers)
- [ ] Timestamps (created_at, updated_at) included where appropriate

### Database & Repository

- [ ] Queries use GORM best practices
- [ ] Proper error handling (check for `gorm.ErrRecordNotFound`)
- [ ] No raw SQL unless necessary
- [ ] Migrations created for schema changes
- [ ] Foreign keys properly defined

### Testing

- [ ] Tests written for new code
- [ ] Repository tests use go-sqlmock
- [ ] Use case tests use mockery mocks
- [ ] Handler tests use httptest
- [ ] Mocks regenerated if interfaces changed (`make mock-gen`)
- [ ] Test coverage acceptable (aim for 80%+)

### Security

- [ ] No sensitive data logged
- [ ] Passwords hashed before storage
- [ ] JWT tokens validated properly
- [ ] SQL injection prevented (using GORM parameterized queries)
- [ ] Input validation on all user inputs
- [ ] No hardcoded secrets or credentials

### Go Best Practices

- [ ] Error handling: errors checked and handled
- [ ] No unused imports or variables
- [ ] Proper variable naming (camelCase for private, PascalCase for public)
- [ ] Comments for exported functions
- [ ] No `panic()` used (except in init/main)
- [ ] Context used for cancellation/timeout where appropriate

## Common Issues & Fixes

### 1. Handler Not Using Unified Response

**Problem**:
```go
c.JSON(200, gin.H{"data": result})
```

**Fix**:
```go
response.NewResponse(c, http.StatusOK, result, "success", nil)
```

### 2. Not Using Language-Specific Errors

**Problem**:
```go
c.JSON(400, gin.H{"error": "user not found"})
```

**Fix**:
```go
lang := middleware.GetLangFromGin(c)
response.NewResponse(c, http.StatusNotFound, nil,
    constants.GetErrorMessage(constants.UserNotFound, lang), nil)
```

### 3. Business Logic in Handler

**Problem**:
```go
func (h *UserHandler) UpdateProfile(c *gin.Context) {
    var req dto.UpdateProfileRequest
    c.ShouldBindJSON(&req)

    // Business logic in handler - BAD!
    user, _ := h.repo.FindByID(userID)
    if user.Email != req.Email {
        // Check email uniqueness
        existing, _ := h.repo.FindByEmail(req.Email)
        if existing != nil {
            c.JSON(400, gin.H{"error": "email taken"})
            return
        }
    }
    user.Email = req.Email
    h.repo.Update(user)
    c.JSON(200, user)
}
```

**Fix**: Move logic to use case
```go
func (h *UserHandler) UpdateProfile(c *gin.Context) {
    lang := middleware.GetLangFromGin(c)

    var req dto.UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.NewResponse(c, http.StatusBadRequest, nil,
            constants.GetErrorMessage(constants.InvalidInput, lang), err.Error())
        return
    }

    claims, _ := c.Get("sess")
    userClaims := claims.(*jwt.JWTClaims)

    // Business logic in use case - GOOD!
    result, err := h.usecase.UpdateProfile(userClaims.UserID, &req)
    if err != nil {
        // Handle errors...
        return
    }

    response.NewResponse(c, http.StatusOK, result, "profile updated", nil)
}
```

### 4. Missing Swagger Comments

**Problem**:
```go
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // No swagger comments
}
```

**Fix**:
```go
// CreateProduct godoc
// @Summary     Create a new product
// @Description Create a new product with the provided details
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       request body dto.CreateProductRequest true "Product details"
// @Success     201 {object} dto.ProductResponse
// @Failure     400 {object} response.Response
// @Router      /products [post]
// @Security    BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // ...
}
```

### 5. Missing DTO Name Comments

**Problem**:
```go
type CreateProductRequest struct {
    Name  string `json:"name"`
    Price float64 `json:"price"`
}
```

**Fix**:
```go
type CreateProductRequest struct {
    Name  string `json:"name" binding:"required"`
    Price float64 `json:"price" binding:"required,min=0"`
} // @name CreateProductRequest
```

### 6. Not Checking Errors

**Problem**:
```go
user, _ := h.repo.FindByID(id)
result, _ := h.usecase.CreateUser(req)
```

**Fix**:
```go
user, err := h.repo.FindByID(id)
if err != nil {
    return nil, err
}

result, err := h.usecase.CreateUser(req)
if err != nil {
    // Handle error appropriately
    return nil, err
}
```

### 7. Hardcoded Values

**Problem**:
```go
const JWTSecret = "my-secret-key"
db := "postgres://user:pass@localhost/db"
```

**Fix**:
```go
// Use config
jwtSecret := config.Load().JWT.Secret
dbURL := config.Load().Database.URL
```

### 8. Wrong Dependency Direction

**Problem**:
```go
// Use case importing infrastructure - BAD!
import "app/internal/shared/infrastructure/repository"

type AuthUsecase struct {
    userRepo *repository.UserRepositoryImpl  // Concrete type!
}
```

**Fix**:
```go
// Use case depending on interface - GOOD!
import "app/internal/shared/domain/repository"

type AuthUsecase struct {
    userRepo repository.UserRepository  // Interface!
}
```

### 9. Missing Validation Tags

**Problem**:
```go
type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Username string `json:"username"`
}
```

**Fix**:
```go
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Username string `json:"username" binding:"required,min=3,max=50"`
}
```

### 10. Not Using Transactions

**Problem**:
```go
// Multiple DB operations without transaction
func (u *OrderUsecase) CreateOrder(req *dto.OrderRequest) error {
    order := &entity.Order{...}
    u.orderRepo.Create(order)

    for _, item := range req.Items {
        u.orderItemRepo.Create(item)  // If this fails, order already created!
    }
    return nil
}
```

**Fix**:
```go
func (u *OrderUsecase) CreateOrder(req *dto.OrderRequest) error {
    return u.db.Transaction(func(tx *gorm.DB) error {
        order := &entity.Order{...}
        if err := tx.Create(order).Error; err != nil {
            return err
        }

        for _, item := range req.Items {
            if err := tx.Create(item).Error; err != nil {
                return err  // Rolls back automatically
            }
        }
        return nil
    })
}
```

## Review Process

### 1. Check Architecture

Read the feature's `module.go` to understand dependencies:

```go
// Good: Dependencies injected, interfaces used
func NewModule(userRepo repository.UserRepository, logger *logrus.Logger) *Module {
    uc := usecase.NewAuthUsecase(userRepo, logger)
    h := handler.NewAuthHandler(uc)
    return &Module{handler: h}
}
```

### 2. Review Handler Layer

Check handlers:
- Use unified response format
- Validate input
- Don't contain business logic
- Have complete Swagger comments
- Handle errors properly

### 3. Review Use Case Layer

Check use cases:
- Contain all business logic
- Depend on repository interfaces
- Return appropriate errors
- Use language-specific error messages
- Are testable

### 4. Review Repository Layer

Check repositories:
- Implement interfaces from domain layer
- Use GORM properly
- Handle errors (especially `gorm.ErrRecordNotFound`)
- Use transactions where needed

### 5. Review Tests

Check tests:
- Cover main scenarios (happy path + error cases)
- Use proper mocks
- Assert expectations
- Have good coverage

### 6. Check Security

- No sensitive data in logs
- Passwords hashed
- JWT properly validated
- Input validated
- No SQL injection risks

## Running Quality Checks

```bash
# Run tests
make test

# Check test coverage
make test-coverage-report

# Format code
go fmt ./...

# Lint code (if golangci-lint installed)
golangci-lint run

# Check for issues
go vet ./...

# Regenerate Swagger after changes
make swag
```

## Approval Criteria

Code is ready to merge when:

1. ✅ Follows Clean Architecture principles
2. ✅ All tests pass with good coverage
3. ✅ Swagger documentation complete
4. ✅ No security issues
5. ✅ Follows project conventions
6. ✅ Error handling comprehensive
7. ✅ Code is readable and maintainable
8. ✅ No hardcoded values
9. ✅ Proper dependency injection
10. ✅ Database migrations created if needed

## Suggesting Improvements

When reviewing, provide:

1. **Specific feedback**: Point to exact lines/files
2. **Explanation**: Why something should change
3. **Example**: Show the correct way
4. **Priority**: Critical vs nice-to-have

**Example feedback format**:
```
In internal/features/auth/delivery/http/handler/auth_handler.go:45

Issue: Business logic in handler
Severity: High
Explanation: Email uniqueness check should be in use case, not handler
Fix: Move this logic to AuthUsecase.Register() method

Before:
[code snippet]

After:
[corrected code snippet]
```
