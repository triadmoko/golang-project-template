package router

import (
	authHandler "app/internal/features/auth/delivery/http/handler"
	"app/internal/features/auth/domain/service"
	productHandler "app/internal/features/product/delivery/http/handler"
	userHandler "app/internal/features/user/delivery/http/handler"
	"app/internal/shared/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router represents the HTTP router
type Router struct {
	authHandler    *authHandler.AuthHandler
	userHandler    *userHandler.UserHandler
	productHandler *productHandler.ProductHandler
	authService    service.AuthService
}

// NewRouter creates a new router
func NewRouter(
	authHandler *authHandler.AuthHandler,
	userHandler *userHandler.UserHandler,
	productHandler *productHandler.ProductHandler,
	authService service.AuthService,
) *Router {
	return &Router{
		authHandler:    authHandler,
		userHandler:    userHandler,
		productHandler: productHandler,
		authService:    authService,
	}
}

// SetupRoutes sets up all the routes
func (r *Router) SetupRoutes() *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin engine
	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(r.authService))
		{
			users.GET("/profile", r.userHandler.GetProfile)
			users.PUT("/profile", r.userHandler.UpdateProfile)
			users.GET("", r.userHandler.GetUsers)
		}

		// Product routes (public)
		products := v1.Group("/products")
		{
			products.GET("", r.productHandler.GetProducts)
			products.GET("/:id", r.productHandler.GetProduct)
			products.GET("/category/:category", r.productHandler.GetProductsByCategory)
			products.GET("/search", r.productHandler.SearchProducts)
			products.POST("", r.productHandler.CreateProduct)
			products.PUT("/:id", r.productHandler.UpdateProduct)
			products.DELETE("/:id", r.productHandler.DeleteProduct)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
