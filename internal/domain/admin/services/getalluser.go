package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// get all users

func GetAllUsers(c *gin.Context) {

	var users []models.User

	if err := database.DB.Where("role = ?", constants.RoleUser).Find(&users).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no users found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": users})
}
