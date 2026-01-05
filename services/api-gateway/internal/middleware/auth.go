package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"monorepo/libs/errors"
	"monorepo/libs/httputil/middleware"
	"monorepo/libs/httputil/response"
	"monorepo/libs/jwt"
)

const (
	UserIDKey       = "user_id"
	UserEmailKey    = "user_email"
	UserUsernameKey = "user_username"
)

// AuthMiddleware creates an authentication middleware using JWT secret
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := middleware.GetLangFromGin(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.NewResponse(c, http.StatusUnauthorized, nil, errors.GetErrorMessage(errors.Unauthorized, lang), nil)
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.NewResponse(c, http.StatusUnauthorized, nil, errors.GetErrorMessage(errors.Unauthorized, lang), nil)
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			response.NewResponse(c, http.StatusUnauthorized, nil, errors.GetErrorMessage(errors.Unauthorized, lang), nil)
			c.Abort()
			return
		}

		// Validate the token
		claims, err := jwt.ValidateToken(secret, token)
		if err != nil {
			response.NewResponse(c, http.StatusUnauthorized, nil, errors.GetErrorMessage(errors.Unauthorized, lang), nil)
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
