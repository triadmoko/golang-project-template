package auth

import (
	"app/internal/features/auth/delivery/http/handler"
	"app/internal/features/auth/usecase"
	"app/internal/shared/domain/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Module is the auth feature module that combines DI and route registration
type Module struct {
	handler *handler.AuthHandler
}

// NewModule creates and wires all auth feature dependencies
func NewModule(userRepo repository.UserRepository, logger *logrus.Logger) *Module {
	// Wire dependencies
	uc := usecase.NewAuthUsecase(userRepo, logger)
	h := handler.NewAuthHandler(uc)

	return &Module{handler: h}
}

// Name returns the feature name
func (m *Module) Name() string {
	return "auth"
}

// RegisterRoutes registers all auth routes
func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	// Auth routes - all public (no auth required)
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", m.handler.Register)
		authGroup.POST("/login", m.handler.Login)
	}
}
