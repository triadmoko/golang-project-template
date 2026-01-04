# Go Project Template

A production-ready Go REST API template with Clean Architecture and feature-based modular design.

## Features

- Clean Architecture with feature-based modules
- JWT authentication
- PostgreSQL with GORM
- Structured logging (Logrus)
- Docker support
- Swagger documentation
- Database migrations
- Multi-language error messages (EN/ID)

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Development Commands](#development-commands)
- [Docker](#docker)
- [Adding New Features](#adding-new-features)
- [Tech Stack](#tech-stack)
- [License](#license)

## Prerequisites

- Go 1.25+
- PostgreSQL 15+
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- Docker & Docker Compose (optional)
- [gow](https://github.com/mitranim/gow) for hot-reload in development (optional)

## Quick Start

```bash
# Clone the repository
git clone <repo-url>
cd golang-project-template

# Setup environment
cp .env.example .env
# Edit .env with your configuration

# Option 1: Docker (recommended)
docker-compose up -d

# Option 2: Manual setup
make migration-up
make dev
```

The API will be available at `http://localhost:8080`.

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `SERVER_HOST` | Server host | `0.0.0.0` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASS` | Database password | `password` |
| `DB_NAME` | Database name | `db_name` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `JWT_SECRET` | JWT signing key | *(required)* |
| `ENV` | Environment | `development` |

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|:----:|-------------|
| `POST` | `/api/v1/auth/register` | No | Register new user |
| `POST` | `/api/v1/auth/login` | No | Login, returns JWT token |
| `GET` | `/api/v1/users/profile` | Yes | Get authenticated user profile |
| `PUT` | `/api/v1/users/profile` | Yes | Update user profile |
| `GET` | `/api/v1/users` | Yes | List users (paginated) |
| `GET` | `/health` | No | Health check |
| `GET` | `/swagger/*` | No | Swagger UI documentation |

**Authentication**: Include JWT token in header: `Authorization: Bearer <token>`

## Project Structure

```
├── cmd/api/                  # Application entry point
├── internal/
│   ├── app/                  # App initialization and routing
│   ├── core/config/          # Configuration management
│   ├── features/             # Feature-based modules
│   │   ├── auth/             # Authentication feature
│   │   │   ├── delivery/     # HTTP handlers & DTOs
│   │   │   └── usecase/      # Business logic
│   │   └── user/             # User management feature
│   │       ├── delivery/     # HTTP handlers & DTOs
│   │       └── usecase/      # Business logic
│   └── shared/               # Shared components
│       ├── domain/           # Entities, repository interfaces, errors
│       ├── infrastructure/   # Database, repository implementations
│       └── delivery/http/    # Middleware, response utilities
├── pkg/                      # Reusable packages
│   ├── jwt/                  # JWT utilities
│   ├── crypto/               # Password hashing
│   └── logger/               # Structured logging
├── migration/                # SQL migration files
└── docs/                     # Swagger documentation
```

## Development Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run with hot-reload (requires gow) |
| `make migration-up` | Run all pending migrations |
| `make migration-down` | Rollback last migration |
| `make migration-create name=xxx` | Create a new migration |
| `make migration-force version=N` | Force migration version |
| `make migration-version` | Show current migration version |
| `make swag` | Generate Swagger documentation |

## Docker

```bash
# Start all services
docker-compose up -d

# View application logs
docker-compose logs -f app

# Stop all services
docker-compose down
```

**Services:**
- `app` - Go application (port 8080)
- `postgres` - PostgreSQL database (port 5432)
- `redis` - Redis cache (port 6379)

## Adding New Features

1. **Create feature directory**
   ```bash
   mkdir -p internal/features/myfeature/{delivery/http/{handler,dto},usecase}
   ```

2. **Implement layers**
   - `domain/` - Entities and repository interfaces (in shared or feature-specific)
   - `usecase/` - Business logic
   - `delivery/http/` - HTTP handlers and DTOs

3. **Create module** with dependency wiring

4. **Register** the module in `internal/app/app.go`

## Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | [Gin](https://github.com/gin-gonic/gin) v1.11 |
| ORM | [GORM](https://gorm.io) v1.25 |
| Database | PostgreSQL |
| Authentication | JWT ([golang-jwt](https://github.com/golang-jwt/jwt) v5) |
| Logging | [Logrus](https://github.com/sirupsen/logrus) |
| Documentation | [Swagger](https://github.com/swaggo/swag) |

## License

MIT License
