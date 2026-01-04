package user

import (
	"app/internal/features/user/delivery/http/handler"
	"app/internal/features/user/usecase"
	"app/internal/shared/delivery/http/middleware"
	"app/internal/shared/domain/repository"

	"github.com/gin-gonic/gin"
)

// Module is the user feature module that combines DI and route registration
type Module struct {
	handler *handler.UserHandler
}

// NewModule creates and wires all user feature dependencies
func NewModule(userRepo repository.UserRepository) *Module {
	// Wire dependencies
	uc := usecase.NewUserUsecase(userRepo)
	h := handler.NewUserHandler(uc)

	return &Module{handler: h}
}

// Name returns the feature name
func (m *Module) Name() string {
	return "user"
}

// RegisterRoutes registers all user routes
func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		// Protected routes - auth middleware applied inline
		users.GET("/profile", middleware.AuthMiddleware(), m.handler.GetProfile)
		users.PUT("/profile", middleware.AuthMiddleware(), m.handler.UpdateProfile)
		users.GET("", middleware.AuthMiddleware(), m.handler.GetUsers)
	}
}
