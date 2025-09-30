package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// user profile

func UserProfile(c *gin.Context) {

	user_id := c.MustGet("user_id").(uint)

	//   find user

	var user models.User

	if err := database.DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": user})
}
