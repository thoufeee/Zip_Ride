package services

import (
	"net/http"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// get all users

func GetAllUsers(c *gin.Context) {

	var users []*models.User

	c.JSON(http.StatusOK, gin.H{"res": users})
}
