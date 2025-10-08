package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// UpdateUser allows an admin to update a user's details
func UpdateUser(c *gin.Context) {
	// Get user ID from param
	userid := c.Param("id")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	// Find the user in the database
	var user models.User
	if err := database.DB.Where("id = ? AND RoleID = ?", userid, constants.RoleUser).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Struct for incoming JSON
	var input struct {
		FirstName   string `json:"firstname"`
		LastName    string `json:"lastname"`
		Email       string `json:"email"`
		Gender      string `json:"gender"`
		PhoneNumber string `json:"phone"`
		Place       string `json:"place"`
		Password    string `json:"password"`
	}

	// Get data to the request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate email format (if updated)
	if input.Email != "" && !utils.EmailCheck(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate phone number (if updated)
	phone, ok := utils.PhoneNumberCheck(input.PhoneNumber)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	// Hash password if changed
	if input.Password != "" {
		hashed, err := utils.GenerateHash(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashed
	}

	// Update the user 
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Gender != "" {
		user.Gender = input.Gender
	}
	if input.PhoneNumber != "" {
		user.PhoneNumber = phone
	}
	if input.Place != "" {
		user.Place = input.Place
	}

	// Save the updated user
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}
