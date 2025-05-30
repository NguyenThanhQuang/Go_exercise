package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: message,
	})
}

func ValidationErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Status:  http.StatusBadRequest,
		Message: "Validation failed",
		Data:    err.Error(),
	})
}
