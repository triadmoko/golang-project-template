package user

import (
	"monorepo/libs/domain/repository"
	"monorepo/services/api-gateway/internal/features/user/delivery/http/handler"
	"monorepo/services/api-gateway/internal/features/user/usecase"
	"monorepo/services/api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Module is the user feature module that combines DI and route registration
type Module struct {
	handler   *handler.UserHandler
	jwtSecret string
}

// NewModule creates and wires all user feature dependencies
func NewModule(userRepo repository.UserRepository, logger *logrus.Logger, jwtSecret string) *Module {
	// Wire dependencies
	uc := usecase.NewUserUsecase(userRepo, logger)
	h := handler.NewUserHandler(uc)

	return &Module{
		handler:   h,
		jwtSecret: jwtSecret,
	}
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
		users.GET("/profile", middleware.AuthMiddleware(m.jwtSecret), m.handler.GetProfile)
		users.PUT("/profile", middleware.AuthMiddleware(m.jwtSecret), m.handler.UpdateProfile)
		users.GET("", middleware.AuthMiddleware(m.jwtSecret), m.handler.GetUsers)
	}
}
