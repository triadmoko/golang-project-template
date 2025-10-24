package middleware

import (
	"app/internal/features/auth/domain/service"
	"app/internal/shared/delivery/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			response.Unauthorized(c, "Token is required")
			c.Abort()
			return
		}

		// Validate the token
		user, err := authService.ValidateToken(token)
		if err != nil {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(UserIDKey, user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_username", user.Username)

		c.Next()
	}
}
