package controllers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

func UpdateStaff(c *gin.Context) {
	staffID := c.Param("id")
	if staffID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Staff id required"})
		return
	}

	var staff models.Admin
	if err := database.DB.Where("id = ? AND RoleID = ?", staffID, constants.RoleAdmin).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	var input struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		PhoneNumber string `json:"phonenumber" binding:"required"`
		Password    string `json:"password"` // optional
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate email
	if !utils.EmailCheck(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Normalize and validate phone number
	phone, ok := utils.PhoneNumberCheck(input.PhoneNumber)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	// Update password if provided
	if input.Password != "" {
		hashed, err := utils.GenerateHash(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		staff.Password = hashed
	}

	// Update fields
	staff.Name = input.Name
	staff.Email = input.Email
	staff.PhoneNumber = phone

	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Staff updated successfully",
		"Name":        staff.Name,
		"Email":       staff.Email,
		"PhoneNumber": staff.PhoneNumber,
	})
}
