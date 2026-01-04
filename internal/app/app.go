package app

import (
	"app/internal/features/auth"
	"app/internal/features/user"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/infrastructure/database"
	sharedRepo "app/internal/shared/infrastructure/repository"
	"app/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Feature defines the interface that each feature module must implement
type Feature interface {
	// Name returns the feature name for logging/debugging
	Name() string
	// RegisterRoutes registers all routes for this feature
	// rg: the base router group (/api/v1)
	RegisterRoutes(rg *gin.RouterGroup)
}

// App holds the application and its dependencies
type App struct {
	DB     *database.PostgresDB
	Engine *gin.Engine
	Logger *logrus.Logger
}

// New creates and initializes the application
func New() (*App, error) {
	app := &App{}

	app.Logger = logger.NewLogger()

	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		return nil, err
	}
	app.DB = db

	// Setup router with features
	app.Engine = app.setupRouter()

	return app, nil
}

// setupRouter configures the HTTP router and registers all features
func (a *App) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LanguageMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Initialize shared repository
	userRepo := sharedRepo.NewUserRepository(a.DB.GetDB())

	// Register all features - just add one line per new feature!
	features := []Feature{
		auth.NewModule(userRepo, a.Logger),
		user.NewModule(userRepo, a.Logger),
	}

	for _, f := range features {
		f.RegisterRoutes(v1)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// Close releases all resources held by the application
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
