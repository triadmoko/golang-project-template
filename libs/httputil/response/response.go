package response

import (
	"github.com/gin-gonic/gin"
)

// Response represents a unified API response
type Response struct {
	Error   bool   `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors,omitempty"`
}

// NewResponse creates a new response with automatic error detection
func NewResponse(c *gin.Context, status int, data any, message string, errs any) {
	isError := status >= 400

	c.JSON(status, Response{
		Error:   isError,
		Status:  status,
		Message: message,
		Data:    data,
		Errors:  errs,
	})
}
