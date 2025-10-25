package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, code int, success bool, message string, data interface{}) {
	c.JSON(code, APIResponse{Success: success, Message: message, Data: data})
}

func Ok(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, true, message, data)
}

func Error(c *gin.Context, code int, message string) {
	JSON(c, code, false, message, nil)
}
