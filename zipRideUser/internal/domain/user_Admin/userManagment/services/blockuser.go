package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// block user

func BlockUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	user.Block = true

	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"res": "user blocked successfuly"})
}

// unblock user
func UnBlockUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	user.Block = false

	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"res": "user unblocked successfuly"})
}
