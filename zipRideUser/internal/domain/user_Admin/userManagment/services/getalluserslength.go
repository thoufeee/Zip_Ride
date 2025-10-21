package services

import (
	"net/http"
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// all useres length

func AllusersLength(c *gin.Context) {

	var user []models.User

	if err := database.DB.Where("role = ?", constants.RoleUser).Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to find all users length"})
		return
	}

	length := len(user)

	c.JSON(http.StatusOK, gin.H{"res": length})
}

// latest users length

func LatestUsersLength(c *gin.Context) {

	var users []models.User

	threeDaysAgo := time.Now().AddDate(0, 0, -3)

	if err := database.DB.Where("role = ? AND created_at >= ?", constants.RoleUser, threeDaysAgo).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to find all users length"})
		return
	}

	length := len(users)

	c.JSON(http.StatusOK, gin.H{"res": length})
}
