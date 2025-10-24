# 🏗️ Feature-Based Project Structure

Project ini telah di-refactor untuk menggunakan struktur berbasis feature yang memisahkan setiap fitur ke dalam direktori terpisah. Ini memungkinkan pengembangan yang lebih modular, mudah di-maintain, dan scalable.

## 📁 Struktur Direktori Lengkap

```
app/
├── cmd/
│   └── api/                    # Application entry point
│       └── main.go
├── internal/
│   ├── features/               # Feature-based modules
│   │   ├── auth/               # 🔐 Authentication feature
│   │   │   ├── domain/
│   │   │   │   ├── entity/     # User entity (untuk auth)
│   │   │   │   ├── repository/ # User repository interface
│   │   │   │   └── service/    # Auth service interface
│   │   │   ├── usecase/        # Auth business logic
│   │   │   ├── infrastructure/
│   │   │   │   ├── repository/ # User repository implementation
│   │   │   │   └── service/    # Auth service implementation
│   │   │   └── delivery/
│   │   │       └── http/
│   │   │           ├── handler/ # Auth HTTP handlers
│   │   │           └── dto/     # Auth DTOs
│   │   ├── user/               # 👤 User management feature
│   │   │   ├── domain/
│   │   │   │   ├── entity/     # User entity (untuk management)
│   │   │   │   └── repository/ # User repository interface
│   │   │   ├── usecase/        # User business logic
│   │   │   ├── infrastructure/
│   │   │   │   └── repository/ # User repository implementation
│   │   │   └── delivery/
│   │   │       └── http/
│   │   │           ├── handler/ # User HTTP handlers
│   │   │           └── dto/     # User DTOs
│   │   └── product/            # 📦 Product management feature
│   │       ├── domain/
│   │       │   ├── entity/     # Product entity
│   │       │   └── repository/ # Product repository interface
│   │       ├── usecase/        # Product business logic
│   │       ├── infrastructure/
│   │       │   └── repository/ # Product repository implementation
│   │       └── delivery/
│   │           └── http/
│   │               ├── handler/ # Product HTTP handlers
│   │               └── dto/     # Product DTOs
│   ├── shared/                 # 🔧 Shared components
│   │   ├── domain/
│   │   │   └── error/          # Shared domain errors
│   │   ├── infrastructure/
│   │   │   └── database/       # Database connection
│   │   └── delivery/
│   │       └── http/
│   │           ├── middleware/ # HTTP middleware
│   │           ├── response/   # Response utilities
│   │           └── router/     # Route configuration
│   └── core/
│       └── config/             # ⚙️ Configuration management
├── pkg/                        # 📦 External packages
│   └── database/               # Database utilities
├── migrations/                 # 🗄️ Database migrations
└── docs/                       # 📚 API documentation
```

## 🎯 Keuntungan Struktur Feature-Based

### ✅ **Modularity**

- Setiap feature terisolasi dengan baik
- Mudah menambah feature baru tanpa mengganggu yang lama
- Tim bisa bekerja pada feature berbeda secara paralel

### ✅ **Maintainability**

- Kode lebih mudah di-maintain karena terorganisir per feature
- Bug fix dan enhancement terbatas pada feature yang relevan
- Dependencies antar feature jelas dan terkontrol

### ✅ **Scalability**

- Mudah dipecah menjadi microservices di masa depan
- Setiap feature bisa di-scale secara independen
- Database per feature bisa dipisah jika diperlukan

### ✅ **Testing**

- Unit test per feature lebih fokus
- Integration test bisa dilakukan per feature
- Mock dependencies lebih mudah dibuat

## 🔧 Cara Menambah Feature Baru

### 1. **Buat Struktur Direktori**

```bash
mkdir -p internal/features/new_feature/{domain/{entity,repository,service},usecase,infrastructure/{repository,service},delivery/http/{handler,dto}}
```

### 2. **Implementasi Layer per Layer**

#### **Domain Layer**

```go
// internal/features/new_feature/domain/entity/entity.go
type NewEntity struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

// internal/features/new_feature/domain/repository/repository.go
type NewRepository interface {
    Create(ctx context.Context, entity *NewEntity) error
    GetByID(ctx context.Context, id string) (*NewEntity, error)
}
```

#### **Use Case Layer**

```go
// internal/features/new_feature/usecase/usecase.go
type NewUsecase interface {
    CreateEntity(ctx context.Context, req *CreateRequest) (*entity.NewEntity, error)
    GetEntity(ctx context.Context, id string) (*entity.NewEntity, error)
}
```

#### **Infrastructure Layer**

```go
// internal/features/new_feature/infrastructure/repository/repository_impl.go
type newRepository struct {
    db *sql.DB
}

func (r *newRepository) Create(ctx context.Context, entity *entity.NewEntity) error {
    // Implementation
}
```

#### **Delivery Layer**

```go
// internal/features/new_feature/delivery/http/handler/handler.go
type NewHandler struct {
    newUsecase usecase.NewUsecase
}

func (h *NewHandler) CreateEntity(c *gin.Context) {
    // Implementation
}

// internal/features/new_feature/delivery/http/dto/dto.go
type CreateRequest struct {
    Name string `json:"name" binding:"required"`
}
```

### 3. **Update Router**

```go
// internal/shared/delivery/http/router/router.go
func (r *Router) SetupRoutes() *gin.Engine {
    // ... existing code ...

    // New feature routes
    newFeature := v1.Group("/new-feature")
    {
        newFeature.POST("", r.newHandler.CreateEntity)
        newFeature.GET("/:id", r.newHandler.GetEntity)
    }

    return router
}
```

### 4. **Update Main**

```go
// cmd/api/main.go
func main() {
    // ... existing code ...

    // Initialize new feature
    newRepo := newRepository.NewRepository(db.DB)
    newUsecase := newUsecase.NewUsecase(newRepo)
    newHandler := newHandler.NewHandler(newUsecase)

    // Add to router
    httpRouter := router.NewRouter(authHandler, userHandler, productHandler, newHandler, authService)
}
```

## 📋 Contoh Feature: Order Management

Berikut contoh implementasi feature Order Management:

### **Domain Layer**

```go
// internal/features/order/domain/entity/order.go
type Order struct {
    ID        uuid.UUID   `json:"id"`
    UserID    uuid.UUID   `json:"user_id"`
    Items     []OrderItem `json:"items"`
    Total     float64     `json:"total"`
    Status    string      `json:"status"`
    CreatedAt time.Time   `json:"created_at"`
}

type OrderItem struct {
    ProductID uuid.UUID `json:"product_id"`
    Quantity  int       `json:"quantity"`
    Price     float64   `json:"price"`
}
```

### **Use Case Layer**

```go
// internal/features/order/usecase/order_usecase.go
type OrderUsecase interface {
    CreateOrder(ctx context.Context, req *CreateOrderRequest) (*entity.Order, error)
    GetOrder(ctx context.Context, orderID string) (*entity.Order, error)
    UpdateOrderStatus(ctx context.Context, orderID string, status string) error
    GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*entity.Order, error)
}
```

### **Delivery Layer**

```go
// internal/features/order/delivery/http/handler/order_handler.go
type OrderHandler struct {
    orderUsecase usecase.OrderUsecase
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var req dto.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request body", err)
        return
    }

    order, err := h.orderUsecase.CreateOrder(c.Request.Context(), &usecase.CreateOrderRequest{
        UserID: req.UserID,
        Items:  req.Items,
    })

    if err != nil {
        // Handle error
        return
    }

    response.Success(c, http.StatusCreated, "Order created successfully", order)
}
```

## 🚀 Best Practices

### 1. **Dependency Direction**

- Domain layer tidak boleh depend ke layer lain
- Use case layer hanya depend ke domain layer
- Infrastructure layer implement interface dari domain layer
- Delivery layer depend ke use case layer

### 2. **Shared Components**

- Gunakan shared components untuk hal yang benar-benar shared
- Hindari circular dependencies antar feature
- Buat interface di shared domain jika diperlukan

### 3. **Naming Convention**

- Package name sesuai dengan feature name
- Handler name: `{Feature}Handler`
- Use case name: `{Feature}Usecase`
- Repository name: `{Entity}Repository`

### 4. **Error Handling**

- Gunakan shared error types untuk consistency
- Custom error per feature jika diperlukan
- Proper error propagation antar layer

## 🔄 Migration dari Struktur Lama

Untuk migrasi dari struktur lama ke feature-based:

1. **Identifikasi Features**: Pisahkan berdasarkan business capability
2. **Move Files**: Pindahkan file ke struktur baru
3. **Update Imports**: Perbaiki semua import paths
4. **Update Dependencies**: Pastikan dependency injection benar
5. **Test**: Pastikan semua test masih berjalan

## 📚 API Endpoints

### **Authentication** (`/api/v1/auth`)

- `POST /register` - Register new user
- `POST /login` - User login

### **Users** (`/api/v1/users`) - Protected

- `GET /profile` - Get user profile
- `PUT /profile` - Update user profile
- `GET /` - Get users list

### **Products** (`/api/v1/products`)

- `GET /` - Get products list
- `GET /:id` - Get product by ID
- `POST /` - Create new product
- `PUT /:id` - Update product
- `DELETE /:id` - Delete product
- `GET /category/:category` - Get products by category
- `GET /search` - Search products

## 🛠️ Development Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run database migrations
make migrate-up

# Generate Swagger documentation
make swagger

# Development setup
make dev-setup
```

## 🐳 Docker Commands

```bash
# Run with Docker Compose
docker-compose up -d

# Build Docker image
make docker-build

# View logs
docker-compose logs -f app
```

## 📖 Referensi

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Modular Monolith](https://martinfowler.com/articles/modular-monolith.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

---

**Struktur ini memungkinkan project untuk berkembang dengan mudah, menambah feature baru tanpa mengganggu yang lama, dan siap untuk dipecah menjadi microservices di masa depan!** 🚀
