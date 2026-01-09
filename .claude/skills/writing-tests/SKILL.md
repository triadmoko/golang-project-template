---
name: writing-tests
description: Write unit tests and integration tests for Go code using testify, go-sqlmock, and mockery. Use when writing tests, creating test files, or testing repositories, use cases, and handlers.
allowed-tools: Read, Write, Edit, Bash, Glob
---

# Writing Tests

This skill guides you through writing tests for this Go project using testify, go-sqlmock, and mockery-generated mocks.

## Test Commands

```bash
make test                      # Run all tests
make test-coverage            # Generate HTML coverage report
make test-coverage-report     # Show coverage in terminal
make mock-gen                 # Generate mocks from interfaces

# Run specific tests
go test ./internal/features/auth/...
go test -run TestCreateUser ./...
go test -v ./internal/shared/infrastructure/repository/
```

## Test File Structure

Test files follow Go conventions:
- Named `*_test.go`
- In the same package as the code being tested
- Use `package packagename` (not `packagename_test`)

```
internal/features/auth/
├── usecase/
│   ├── auth_usecase.go
│   └── auth_usecase_test.go      # Test file
└── delivery/http/handler/
    ├── auth_handler.go
    └── auth_handler_test.go       # Test file
```

## Repository Tests (with go-sqlmock)

Repository tests use go-sqlmock to mock database interactions without a real database.

### Example: Testing a Repository

```go
package repository

import (
    "app/internal/shared/domain/entity"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
    // Create mock DB
    mockDB, mock, err := sqlmock.New()
    assert.NoError(t, err)

    // Create GORM DB with mock
    dialector := postgres.New(postgres.Config{
        Conn:       mockDB,
        DriverName: "postgres",
    })

    db, err := gorm.Open(dialector, &gorm.Config{})
    assert.NoError(t, err)

    return db, mock
}

func TestUserRepository_Create(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepository(db)

    user := &entity.User{
        ID:       uuid.New(),
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
        FullName: "Test User",
    }

    // Expect INSERT query
    mock.ExpectBegin()
    mock.ExpectExec(`INSERT INTO "users"`).
        WithArgs(
            user.ID,
            user.Username,
            user.Email,
            user.Password,
            user.FullName,
            sqlmock.AnyArg(), // created_at
            sqlmock.AnyArg(), // updated_at
        ).
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    // Execute
    err := repo.Create(user)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepository(db)

    expectedUser := &entity.User{
        ID:       uuid.New(),
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
        FullName: "Test User",
    }

    // Define expected query result
    rows := sqlmock.NewRows([]string{
        "id", "username", "email", "password", "full_name", "created_at", "updated_at",
    }).AddRow(
        expectedUser.ID,
        expectedUser.Username,
        expectedUser.Email,
        expectedUser.Password,
        expectedUser.FullName,
        time.Now(),
        time.Now(),
    )

    // Expect SELECT query
    mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
        WithArgs("test@example.com").
        WillReturnRows(rows)

    // Execute
    result, err := repo.FindByEmail("test@example.com")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedUser.Email, result.Email)
    assert.Equal(t, expectedUser.Username, result.Username)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepository(db)

    user := &entity.User{
        ID:       uuid.New(),
        Username: "testuser",
        Email:    "test@example.com",
        FullName: "Updated Name",
    }

    // Expect UPDATE query
    mock.ExpectBegin()
    mock.ExpectExec(`UPDATE "users" SET`).
        WithArgs(
            user.Username,
            user.Email,
            user.FullName,
            sqlmock.AnyArg(), // updated_at
            user.ID,
        ).
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    // Execute
    err := repo.Update(user)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewUserRepository(db)

    userID := uuid.New()

    // Expect DELETE query
    mock.ExpectBegin()
    mock.ExpectExec(`DELETE FROM "users" WHERE "users"."id" = \$1`).
        WithArgs(userID).
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    // Execute
    err := repo.Delete(userID)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### Common sqlmock Patterns

**Expect query with specific columns**:
```go
rows := sqlmock.NewRows([]string{"id", "name", "email"}).
    AddRow(1, "John", "john@example.com").
    AddRow(2, "Jane", "jane@example.com")

mock.ExpectQuery(`SELECT \* FROM users`).WillReturnRows(rows)
```

**Expect error**:
```go
mock.ExpectQuery(`SELECT \* FROM users`).
    WillReturnError(gorm.ErrRecordNotFound)
```

**Transaction expectations**:
```go
mock.ExpectBegin()
mock.ExpectExec(`INSERT INTO...`).WillReturnResult(sqlmock.NewResult(1, 1))
mock.ExpectCommit()
```

## Use Case Tests (with Mocks)

Use case tests use mockery-generated mocks to test business logic without dependencies.

### Step 1: Generate Mocks

```bash
make mock-gen
```

This generates mocks in `internal/mocks/` based on `.mockery.yaml`.

### Step 2: Write Use Case Tests

```go
package usecase

import (
    "app/internal/features/auth/delivery/http/dto"
    "app/internal/shared/constants"
    "app/internal/shared/domain/entity"
    mocks "app/internal/mocks/repository"
    "errors"
    "testing"

    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestAuthUsecase_Register_Success(t *testing.T) {
    // Setup
    mockRepo := mocks.NewMockUserRepository(t)
    logger := logrus.New()
    usecase := NewAuthUsecase(mockRepo, logger)

    req := &dto.RegisterRequest{
        Username: "newuser",
        Email:    "new@example.com",
        Password: "password123",
        FullName: "New User",
    }

    // Mock expectations
    mockRepo.EXPECT().
        FindByEmail(req.Email).
        Return(nil, gorm.ErrRecordNotFound).
        Once()

    mockRepo.EXPECT().
        FindByUsername(req.Username).
        Return(nil, gorm.ErrRecordNotFound).
        Once()

    mockRepo.EXPECT().
        Create(mock.AnythingOfType("*entity.User")).
        Return(nil).
        Once()

    // Execute
    result, err := usecase.Register(req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, req.Email, result.Email)
    assert.Equal(t, req.Username, result.Username)
    mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register_EmailAlreadyExists(t *testing.T) {
    // Setup
    mockRepo := mocks.NewMockUserRepository(t)
    logger := logrus.New()
    usecase := NewAuthUsecase(mockRepo, logger)

    req := &dto.RegisterRequest{
        Username: "newuser",
        Email:    "existing@example.com",
        Password: "password123",
    }

    existingUser := &entity.User{
        ID:    uuid.New(),
        Email: req.Email,
    }

    // Mock expectations - email already exists
    mockRepo.EXPECT().
        FindByEmail(req.Email).
        Return(existingUser, nil).
        Once()

    // Execute
    result, err := usecase.Register(req)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Equal(t, constants.GetError(constants.UserAlreadyExists, constants.LangEN), err)
    mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Success(t *testing.T) {
    // Setup
    mockRepo := mocks.NewMockUserRepository(t)
    logger := logrus.New()
    usecase := NewAuthUsecase(mockRepo, logger)

    hashedPassword, _ := crypto.HashPassword("password123")
    user := &entity.User{
        ID:       uuid.New(),
        Email:    "user@example.com",
        Password: hashedPassword,
    }

    req := &dto.LoginRequest{
        Email:    "user@example.com",
        Password: "password123",
    }

    // Mock expectations
    mockRepo.EXPECT().
        FindByEmail(req.Email).
        Return(user, nil).
        Once()

    // Execute
    result, err := usecase.Login(req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.NotEmpty(t, result.Token)
    mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_InvalidPassword(t *testing.T) {
    // Setup
    mockRepo := mocks.NewMockUserRepository(t)
    logger := logrus.New()
    usecase := NewAuthUsecase(mockRepo, logger)

    hashedPassword, _ := crypto.HashPassword("correctpassword")
    user := &entity.User{
        ID:       uuid.New(),
        Email:    "user@example.com",
        Password: hashedPassword,
    }

    req := &dto.LoginRequest{
        Email:    "user@example.com",
        Password: "wrongpassword",
    }

    // Mock expectations
    mockRepo.EXPECT().
        FindByEmail(req.Email).
        Return(user, nil).
        Once()

    // Execute
    result, err := usecase.Login(req)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    mockRepo.AssertExpectations(t)
}
```

### Mock Expectations Patterns

**Basic expectation**:
```go
mockRepo.EXPECT().
    MethodName(arg1, arg2).
    Return(result, nil).
    Once()
```

**Any argument**:
```go
mockRepo.EXPECT().
    Create(mock.AnythingOfType("*entity.User")).
    Return(nil).
    Once()
```

**Multiple calls**:
```go
mockRepo.EXPECT().FindByID(userID).Return(user, nil).Times(3)
```

**Different returns per call**:
```go
mockRepo.EXPECT().FindByID(userID).Return(user, nil).Once()
mockRepo.EXPECT().FindByID(userID).Return(nil, errors.New("error")).Once()
```

**Argument matcher**:
```go
mockRepo.EXPECT().
    Create(mock.MatchedBy(func(u *entity.User) bool {
        return u.Email == "test@example.com"
    })).
    Return(nil).
    Once()
```

## Handler Tests (HTTP Tests)

Handler tests use `httptest` to test HTTP endpoints without a running server.

```go
package handler

import (
    "app/internal/features/auth/delivery/http/dto"
    mocks "app/internal/mocks/usecase"
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    return gin.New()
}

func TestAuthHandler_Register_Success(t *testing.T) {
    // Setup
    mockUsecase := mocks.NewMockAuthUsecase(t)
    handler := NewAuthHandler(mockUsecase)
    router := setupRouter()
    router.POST("/register", handler.Register)

    req := dto.RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
        FullName: "Test User",
    }

    expectedResponse := &dto.RegisterResponse{
        ID:       "123e4567-e89b-12d3-a456-426614174000",
        Username: req.Username,
        Email:    req.Email,
    }

    // Mock expectations
    mockUsecase.EXPECT().
        Register(&req).
        Return(expectedResponse, nil).
        Once()

    // Create request
    body, _ := json.Marshal(req)
    w := httptest.NewRecorder()
    r := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
    r.Header.Set("Content-Type", "application/json")

    // Execute
    router.ServeHTTP(w, r)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.False(t, response["error"].(bool))
    assert.NotNil(t, response["data"])

    mockUsecase.AssertExpectations(t)
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
    // Setup
    mockUsecase := mocks.NewMockAuthUsecase(t)
    handler := NewAuthHandler(mockUsecase)
    router := setupRouter()
    router.POST("/register", handler.Register)

    // Invalid request (missing required fields)
    req := dto.RegisterRequest{
        Username: "", // Empty username
        Email:    "invalid-email", // Invalid email
    }

    // Create request
    body, _ := json.Marshal(req)
    w := httptest.NewRecorder()
    r := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
    r.Header.Set("Content-Type", "application/json")

    // Execute
    router.ServeHTTP(w, r)

    // Assert
    assert.Equal(t, http.StatusBadRequest, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.True(t, response["error"].(bool))
    assert.NotNil(t, response["errors"])
}

func TestAuthHandler_Login_Success(t *testing.T) {
    // Setup
    mockUsecase := mocks.NewMockAuthUsecase(t)
    handler := NewAuthHandler(mockUsecase)
    router := setupRouter()
    router.POST("/login", handler.Login)

    req := dto.LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }

    expectedResponse := &dto.LoginResponse{
        Token: "jwt.token.here",
        User: dto.UserInfo{
            ID:       "123e4567-e89b-12d3-a456-426614174000",
            Email:    req.Email,
            Username: "testuser",
        },
    }

    // Mock expectations
    mockUsecase.EXPECT().
        Login(&req).
        Return(expectedResponse, nil).
        Once()

    // Create request
    body, _ := json.Marshal(req)
    w := httptest.NewRecorder()
    r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
    r.Header.Set("Content-Type", "application/json")

    // Execute
    router.ServeHTTP(w, r)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.False(t, response["error"].(bool))

    mockUsecase.AssertExpectations(t)
}
```

## Test Organization

### Table-Driven Tests

For testing multiple scenarios:

```go
func TestUserRepository_FindByEmail(t *testing.T) {
    tests := []struct {
        name          string
        email         string
        setupMock     func(sqlmock.Sqlmock)
        expectedError bool
    }{
        {
            name:  "user found",
            email: "found@example.com",
            setupMock: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"id", "email"}).
                    AddRow(uuid.New(), "found@example.com")
                mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
            },
            expectedError: false,
        },
        {
            name:  "user not found",
            email: "notfound@example.com",
            setupMock: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)
            },
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, mock := setupMockDB(t)
            repo := NewUserRepository(db)

            tt.setupMock(mock)

            result, err := repo.FindByEmail(tt.email)

            if tt.expectedError {
                assert.Error(t, err)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
            }
        })
    }
}
```

## Running and Analyzing Tests

### Run tests with coverage
```bash
make test-coverage
```

### View coverage in browser
```bash
make test-coverage
open coverage.html
```

### Run specific package
```bash
go test -v ./internal/features/auth/usecase/
```

### Run specific test
```bash
go test -v -run TestAuthUsecase_Register_Success ./internal/features/auth/usecase/
```

### Run with race detector
```bash
go test -race ./...
```

## Best Practices

1. **Test file naming**: Use `*_test.go` suffix
2. **Test function naming**: `Test{FunctionName}_{Scenario}`
3. **Mock generation**: Run `make mock-gen` after changing interfaces
4. **Coverage target**: Aim for 80%+ coverage
5. **Arrange-Act-Assert**: Structure tests clearly
6. **Test one thing**: Each test should verify one behavior
7. **Use table-driven tests**: For multiple similar scenarios
8. **Clean up mocks**: Always call `AssertExpectations(t)`
