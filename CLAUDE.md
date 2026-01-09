# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Development
make dev                          # Run with hot-reload using gow
make test                         # Run all tests
make test-coverage                # Generate coverage report as HTML
make test-coverage-report         # Display coverage in terminal

# Database Migrations
make migration-up                 # Run all pending migrations
make migration-down               # Rollback last migration
make migration-create name=xxx    # Create new migration file
make migration-force version=N    # Force migration to specific version
make migration-version            # Show current migration version

# Documentation
make swag                         # Regenerate Swagger docs (required after changing API comments)

# Mocks
make mock-gen                     # Generate mocks from .mockery.yaml
make mock-clean                   # Remove all generated mocks
```

## Architecture Overview

This codebase follows **Clean Architecture** with a **feature-based modular design**. Each feature is self-contained with its own handlers, DTOs, and use cases, while sharing domain entities and repositories.

### Dependency Flow

```
cmd/api/main.go
  └─> internal/app/app.go (initializes all features)
       └─> Feature Modules (auth, user)
            └─> Handlers (HTTP layer)
                 └─> Use Cases (business logic)
                      └─> Repository Interfaces (domain layer)
                           └─> Repository Implementations (infrastructure)
```

**Key principle**: Dependencies point inward. Use cases depend on repository interfaces (not implementations), and handlers depend on use case interfaces.

### Feature Module Pattern

Each feature follows a standard structure and implements the `Feature` interface defined in `internal/app/app.go`:

```go
type Feature interface {
    Name() string
    RegisterRoutes(rg *gin.RouterGroup)
}
```

A feature module (`module.go`) handles dependency injection and route registration. Example structure:

```
internal/features/myfeature/
├── module.go              # DI container + route registration
├── delivery/http/
│   ├── handler/           # HTTP handlers
│   └── dto/               # Request/response DTOs
└── usecase/               # Business logic + interface definition
```

**To add a new feature:**
1. Create the feature directory structure
2. Implement entities/repositories in `internal/shared/domain/` if needed
3. Create use case interface and implementation
4. Create handler and DTOs
5. Create `module.go` with DI wiring
6. Register in `internal/app/app.go` features slice

See `internal/features/auth/module.go` for reference.

### Shared Components

**Domain Layer** (`internal/shared/domain/`):
- `entity/` - Domain entities (User, etc.)
- `repository/` - Repository interfaces
- `errors/` - Domain-specific errors

**Infrastructure** (`internal/shared/infrastructure/`):
- `database/` - Database connection setup
- `repository/` - Repository implementations (depends on domain interfaces)

**HTTP Utilities** (`internal/shared/delivery/http/`):
- `middleware/` - Authentication, CORS, logging, language detection
- `response/` - Unified API response format

**Reusable Packages** (`pkg/`):
- `jwt/` - Token generation and validation
- `crypto/` - Password hashing (bcrypt)
- `logger/` - Structured logging setup
- `pagination.go` - Pagination utilities

### Authentication Flow

1. JWT tokens are issued via `/api/v1/auth/login`
2. Protected routes use `middleware.AuthMiddleware()`
3. Middleware validates token and stores claims in context under key `"sess"`
4. Handlers retrieve claims via `c.Get("sess")` (returns `*jwt.JWTClaims`)

### Multi-Language Error Messages

Error messages support EN/ID via `internal/shared/constants/error.go`:

```go
constants.GetError(constants.UserNotFound, lang)        // returns error
constants.GetErrorMessage(constants.Unauthorized, lang) // returns string
```

Language is detected from `Accept-Language` header by `middleware.LanguageMiddleware()`.

### Testing

- Mock generation uses mockery (config: `.mockery.yaml`)
- Run `make mock-gen` after changing interfaces
- Test files use go-sqlmock for database mocking
- See `internal/shared/infrastructure/repository/user_repository_impl_test.go` for examples

### Swagger Documentation

- Swagger annotations are in `cmd/api/main.go` (global) and handler methods
- Run `make swag` to regenerate docs after API changes
- Access at `http://localhost:8080/swagger/index.html`
- Use `@Security BearerAuth` annotation for protected endpoints

### Configuration

All config loaded via `internal/core/config/` from environment variables. Required variables:
- `DB_*` - Database connection
- `JWT_SECRET` - Token signing key
- `SERVER_PORT` / `SERVER_HOST` - Server binding

## Development Notes

- The module path is `app` (see go.mod)
- Database uses GORM with PostgreSQL
- Gin runs in release mode by default
- Server implements graceful shutdown
- All routes are versioned under `/api/v1`
