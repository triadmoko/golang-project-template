package main

import (
	"app/internal/core/config"
	authHandler "app/internal/features/auth/delivery/http/handler"
	authRepository "app/internal/features/auth/infrastructure/repository"
	authService "app/internal/features/auth/infrastructure/service"
	authUsecase "app/internal/features/auth/usecase"
	productHandler "app/internal/features/product/delivery/http/handler"
	productRepository "app/internal/features/product/infrastructure/repository"
	productUsecase "app/internal/features/product/usecase"
	userHandler "app/internal/features/user/delivery/http/handler"
	userRepository "app/internal/features/user/infrastructure/repository"
	userUsecase "app/internal/features/user/usecase"
	"app/internal/shared/delivery/http/router"
	"app/internal/shared/infrastructure/database"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title Yopatungan Backend API
// @version 1.0
// @description A REST API built with Go and Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	authUserRepo := authRepository.NewUserRepository(db.GetDB())
	userRepo := userRepository.NewUserRepository(db.GetDB())
	productRepo := productRepository.NewProductRepository(db.GetDB())

	// Initialize services
	authService := authService.NewAuthService(cfg.JWT.Secret)

	// Initialize use cases
	authUsecase := authUsecase.NewAuthUsecase(authUserRepo, authService)
	userUsecase := userUsecase.NewUserUsecase(userRepo)
	productUsecase := productUsecase.NewProductUsecase(productRepo)

	// Initialize handlers
	authHandler := authHandler.NewAuthHandler(authUsecase)
	userHandler := userHandler.NewUserHandler(userUsecase)
	productHandler := productHandler.NewProductHandler(productUsecase)

	// Initialize router
	httpRouter := router.NewRouter(authHandler, userHandler, productHandler, authService)
	ginEngine := httpRouter.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: ginEngine,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
