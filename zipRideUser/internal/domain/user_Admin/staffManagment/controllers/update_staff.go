package controllers

import (
	"encoding/json"
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// update admin profile
func UpdateStaff(c *gin.Context) {
	staffID := c.Param("id")
	if staffID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Staff id required"})
		return
	}

	var staff models.Admin
<<<<<<< HEAD
	if err := database.DB.Where("id = ? AND role = ?", staffID, constants.RoleAdmin).First(&staff).Error; err != nil {
=======
	if err := database.DB.Where("id = ? AND RoleID = ?", staffID, constants.RoleAdmin).First(&staff).Error; err != nil {
>>>>>>> 2c00f30 (folders changed)
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	var input struct {
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		PhoneNumber string   `json:"phonenumber"`
		Password    string   `json:"password"`
		ExtraPerms  []string `json:"extra_permissions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate email
	if input.Email != "" && !utils.EmailCheck(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if input.PhoneNumber != "" {
		phone, ok := utils.PhoneNumberCheck(input.PhoneNumber)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
			return
		}
		staff.PhoneNumber = phone
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

	if input.Name != "" {
		staff.Name = input.Name
	}

	if len(input.ExtraPerms) > 0 {
		var existing []string

		if err := json.Unmarshal([]byte(staff.Permissions), &existing); err != nil {
			existing = []string{}
		}

		perMap := make(map[string]bool)
		for _, p := range existing {
			{
				perMap[p] = true
			}
		}

		for _, p := range input.ExtraPerms {
			if p != "Click an available permission to add it." {
				perMap[p] = true
			}
		}

		merged := []string{}
		for p := range perMap {
			merged = append(merged, p)
		}

		permjson, err := json.Marshal(merged)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to process permissions"})
			return
		}

		staff.Permissions = permjson
	}

	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Staff updated successfully",
		"Name":    staff.Name,
		"Email":   staff.Email,
	})
}
