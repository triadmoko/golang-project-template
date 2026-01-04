package middleware

import (
	"app/internal/core/config"
	"app/internal/shared/delivery/http/response"
	"app/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey       = "user_id"
	UserEmailKey    = "user_email"
	UserUsernameKey = "user_username"
)

// AuthMiddleware creates an authentication middleware using JWT secret
func AuthMiddleware() gin.HandlerFunc {
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
		claims, err := jwt.ValidateToken(config.Load().JWT.Secret, token)
		if err != nil {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// Set user information in context (all strings now)
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserUsernameKey, claims.Username)

		c.Next()
	}
}
