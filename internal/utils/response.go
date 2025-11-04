package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Code:    code,
	})
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

