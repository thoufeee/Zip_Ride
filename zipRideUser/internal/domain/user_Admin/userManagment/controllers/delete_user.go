package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// DeleteUser allows admin to delete a user by ID
func DeleteUser(c *gin.Context) {
	// Get user ID from param
	userid := c.Param("id")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	var user models.User

	// Check if user exists in database
	if err := database.DB.Where("id = ?", userid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete user
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"user":    user,
	})
}
