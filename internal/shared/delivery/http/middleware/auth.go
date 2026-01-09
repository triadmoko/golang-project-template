package middleware

import (
	"app/internal/core/config"
	"app/internal/shared/constants"
	"app/internal/shared/delivery/http/response"
	"app/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	SESS = "sess"
)

// AuthMiddleware creates an authentication middleware using JWT secret
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := GetLangFromGin(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
			c.Abort()
			return
		}

		// Validate the token
		claims, err := jwt.ValidateToken(config.Load().JWT.Secret, token)
		if err != nil {
			response.NewResponse(c, http.StatusUnauthorized, nil, constants.GetErrorMessage(constants.Unauthorized, lang), nil)
			c.Abort()
			return
		}

		// Set user information in context (all strings now)
		c.Set(SESS, claims)
		c.Next()
	}
}
