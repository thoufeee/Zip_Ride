package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/middleware"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// delete user profile

func DeleteUserProfile(c *gin.Context) {

	user_id := middleware.GetUserID(c)

	if user_id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unathorized"})
		return
	}

	var user models.User

	if err := database.DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to delete user account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly deleted user account"})
}
