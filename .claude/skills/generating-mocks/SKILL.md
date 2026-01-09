---
name: generating-mocks
description: Generate mock implementations using mockery for testing. Use when creating mocks, updating mock interfaces, or setting up test doubles.
allowed-tools: Read, Write, Edit, Bash
---

# Generating Mocks

This skill guides you through generating and managing mock implementations using mockery.

## Quick Commands

```bash
make mock-gen      # Generate all mocks
make mock-clean    # Remove all generated mocks
```

## Mockery Configuration

Configuration is in `.mockery.yaml`:

```yaml
with-expecter: true    # Enable EXPECT() style assertions
testonly: false        # Don't add //go:build test constraint
all: false             # Don't auto-generate for all interfaces
packages:
  app/internal/shared/domain/repository:
    interfaces:
      UserRepository:
        config:
          dir: internal/mocks/repository
          outpkg: mocks
  app/internal/features/auth/usecase:
    interfaces:
      AuthUsecase:
        config:
          dir: internal/mocks/usecase
          outpkg: mocks
```

## Adding a New Interface to Mock

### Step 1: Define the Interface

Create your interface in the appropriate package:

```go
// internal/features/product/usecase/product_usecase.go
package usecase

import "app/internal/features/product/delivery/http/dto"

type ProductUsecase interface {
    CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
    GetProductByID(id string) (*dto.ProductResponse, error)
    ListProducts(page, perPage int) ([]dto.ProductResponse, int64, error)
    UpdateProduct(id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
    DeleteProduct(id string) error
}
```

### Step 2: Add to .mockery.yaml

Add the interface to `.mockery.yaml`:

```yaml
packages:
  # ... existing packages ...

  app/internal/features/product/usecase:
    interfaces:
      ProductUsecase:
        config:
          dir: internal/mocks/usecase
          outpkg: mocks
```

### Step 3: Generate Mocks

```bash
make mock-gen
```

This creates: `internal/mocks/usecase/mock_ProductUsecase.go`

### Step 4: Use the Mock in Tests

```go
package handler

import (
    mocks "app/internal/mocks/usecase"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestProductHandler_CreateProduct(t *testing.T) {
    // Create mock
    mockUsecase := mocks.NewMockProductUsecase(t)

    // Set expectations
    mockUsecase.EXPECT().
        CreateProduct(mock.AnythingOfType("*dto.CreateProductRequest")).
        Return(&dto.ProductResponse{ID: "123"}, nil).
        Once()

    // Use mock in tests...
    handler := NewProductHandler(mockUsecase)

    // ... test code ...

    // Verify expectations
    mockUsecase.AssertExpectations(t)
}
```

## Mock Patterns

### Repository Mocks

Repository interfaces are typically in `internal/shared/domain/repository/`:

```yaml
packages:
  app/internal/shared/domain/repository:
    interfaces:
      ProductRepository:
        config:
          dir: internal/mocks/repository
          outpkg: mocks
```

Generated mock: `internal/mocks/repository/mock_ProductRepository.go`

### Use Case Mocks

Use case interfaces are in feature-specific packages:

```yaml
packages:
  app/internal/features/product/usecase:
    interfaces:
      ProductUsecase:
        config:
          dir: internal/mocks/usecase
          outpkg: mocks
```

Generated mock: `internal/mocks/usecase/mock_ProductUsecase.go`

## Using Generated Mocks

### Basic Expectations

```go
// Exact arguments
mockRepo.EXPECT().
    FindByID("123").
    Return(user, nil).
    Once()

// Any arguments
mockRepo.EXPECT().
    Create(mock.AnythingOfType("*entity.User")).
    Return(nil).
    Once()

// Multiple calls
mockRepo.EXPECT().
    List().
    Return([]entity.User{}, nil).
    Times(3)
```

### Argument Matchers

```go
import "github.com/stretchr/testify/mock"

// Match any value of type
mock.AnythingOfType("string")
mock.AnythingOfType("*entity.User")

// Custom matcher
mockRepo.EXPECT().
    Create(mock.MatchedBy(func(u *entity.User) bool {
        return u.Email == "test@example.com"
    })).
    Return(nil).
    Once()

// Anything (any type)
mockRepo.EXPECT().
    Update(mock.Anything).
    Return(nil).
    Once()
```

### Return Values

```go
// Return single value
mockRepo.EXPECT().
    Delete(userID).
    Return(nil).
    Once()

// Return multiple values
mockRepo.EXPECT().
    FindByEmail("test@example.com").
    Return(user, nil).
    Once()

// Return error
mockRepo.EXPECT().
    FindByID("invalid").
    Return(nil, gorm.ErrRecordNotFound).
    Once()

// Different returns for sequential calls
mockRepo.EXPECT().FindAll().Return([]User{user1}, nil).Once()
mockRepo.EXPECT().FindAll().Return([]User{user1, user2}, nil).Once()
```

### Call Counts

```go
// Called exactly once (default)
mockRepo.EXPECT().Create(user).Return(nil).Once()

// Called exactly N times
mockRepo.EXPECT().Update(user).Return(nil).Times(3)

// Called any number of times (including 0)
mockRepo.EXPECT().GetConfig().Return("config").Maybe()
```

## Complete Example

### 1. Define Interface

```go
// internal/features/notification/usecase/notification_usecase.go
package usecase

type NotificationUsecase interface {
    SendEmail(to, subject, body string) error
    SendSMS(to, message string) error
}
```

### 2. Add to .mockery.yaml

```yaml
packages:
  app/internal/features/notification/usecase:
    interfaces:
      NotificationUsecase:
        config:
          dir: internal/mocks/usecase
          outpkg: mocks
```

### 3. Generate Mock

```bash
make mock-gen
```

### 4. Use in Test

```go
package usecase

import (
    "app/internal/features/notification/usecase"
    mocks "app/internal/mocks/usecase"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserUsecase_RegisterWithNotification(t *testing.T) {
    // Setup mocks
    mockUserRepo := mocks.NewMockUserRepository(t)
    mockNotification := mocks.NewMockNotificationUsecase(t)

    userUsecase := NewUserUsecase(mockUserRepo, mockNotification)

    // Set expectations
    mockUserRepo.EXPECT().
        Create(mock.AnythingOfType("*entity.User")).
        Return(nil).
        Once()

    mockNotification.EXPECT().
        SendEmail(
            "newuser@example.com",
            "Welcome!",
            mock.AnythingOfType("string"),
        ).
        Return(nil).
        Once()

    // Execute test
    err := userUsecase.Register(&RegisterRequest{
        Email: "newuser@example.com",
    })

    // Assert
    assert.NoError(t, err)
    mockUserRepo.AssertExpectations(t)
    mockNotification.AssertExpectations(t)
}
```

## Regenerating Mocks

### When to Regenerate

Run `make mock-gen` after:
- Adding new methods to an interface
- Changing method signatures
- Removing methods from an interface
- Adding new interfaces to mock

### Handling Changes

1. **Interface method added**: Regenerate mocks, update tests
2. **Interface method removed**: Regenerate mocks, remove test expectations
3. **Method signature changed**: Regenerate mocks, update test expectations

```bash
# Clean and regenerate all mocks
make mock-clean
make mock-gen
```

## Troubleshooting

### Mock Not Found

**Problem**: `undefined: mocks.NewMockUserRepository`

**Solution**:
1. Check interface is in `.mockery.yaml`
2. Run `make mock-gen`
3. Check mock file exists in `internal/mocks/`

### Wrong Package Name

**Problem**: Mock generated with wrong package

**Solution**: Check `outpkg` in `.mockery.yaml`:
```yaml
config:
  dir: internal/mocks/repository
  outpkg: mocks  # Must be "mocks"
```

### Interface Not Found

**Problem**: `mockery: error: failed to find interface UserRepository`

**Solution**:
1. Verify interface exists and is exported (starts with capital letter)
2. Check package path in `.mockery.yaml` matches actual package
3. Ensure interface name is spelled correctly

### Expectation Not Met

**Problem**: Test fails with "FAIL: Expected function call but did not find it"

**Solution**:
1. Ensure method is actually called in code
2. Check argument matchers match actual arguments
3. Use `mock.Anything` for flexible matching during debugging
4. Add `.Maybe()` if call is optional

## Best Practices

1. **Always call `AssertExpectations(t)`** at the end of tests
2. **Use specific matchers** when possible for better test clarity
3. **Mock only what you need** - don't set expectations for unused methods
4. **Regenerate after interface changes** to keep mocks in sync
5. **Use `Once()`** explicitly to catch unexpected multiple calls
6. **Group related mocks** in the same directory structure
7. **Keep `.mockery.yaml` organized** with clear package groupings

## Mock Directory Structure

```
internal/mocks/
├── repository/
│   ├── mock_UserRepository.go
│   ├── mock_ProductRepository.go
│   └── mock_OrderRepository.go
└── usecase/
    ├── mock_AuthUsecase.go
    ├── mock_UserUsecase.go
    └── mock_ProductUsecase.go
```

This structure keeps mocks organized by layer (repository vs usecase) for easy discovery.
