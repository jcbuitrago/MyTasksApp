package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendResponse sends a standardized JSON response
func SendResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	})
}

// SendError sends a standardized error response
func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Status:  http.StatusText(statusCode),
		Message: message,
	})
}