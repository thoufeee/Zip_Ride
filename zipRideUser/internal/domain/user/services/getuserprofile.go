package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/middleware"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// getting user profile

func GetUserProfile(c *gin.Context) {

	userId := middleware.GetUserID(c)

	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unathorized"})
		return
	}

	var user models.User

	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": user})
}
