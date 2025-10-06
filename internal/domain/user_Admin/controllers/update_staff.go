package controllers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/middleware"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

func UpdateStaff(c *gin.Context) {
	//Get staff if from param
	staffid := middleware.GetUserID(c)
	//get staff data from the database
	var staff models.Admin
	if err := database.DB.Where("id = ? AND RoleID = ?", staffid, constants.RoleStaff).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	//staff update input struct
	var input struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		PhoneNumber string `json:"phonenumber" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}
	//get the input data to update staff details
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inputs"})
		return
	}
	//check email to veryfy if it is emty and email check
	if input.Email != "" && !utils.EmailCheck(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}
	//check phone number to verify if it is empty and phone number check function
	if input.PhoneNumber != "" && !utils.PhoneNumberCheck(input.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number format"})
		return
	}
	// password if it is changing if it is empty else to generate hashing to secure
	if input.Password != "" {
		hashed, err := utils.GenerateHash(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		staff.Password = hashed
	}
	//update field provided
	//update name
	if input.Name != "" {
		staff.Name = input.Name
	}
	// update email 
	if input.Email != "" {
		staff.Email = input.Email
	}
	// update phone number 
	if input.PhoneNumber != "" {
		staff.PhoneNumber = input.PhoneNumber
	}

	//save to database
	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff"})
		return
	}

	//sucess Responce
	c.JSON(http.StatusOK, gin.H{
		"message": "Staff update sucessfully",
		"staff": staff,
	})
}
