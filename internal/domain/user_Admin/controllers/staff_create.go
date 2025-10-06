package controllers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// creating staff || manager
func CreateStaff(c *gin.Context) {

		var data struct {
			Name        string `json:"name" binding:"required"`
			Email       string `json:"email" binding:"required"`
			PhoneNumber string `json:"phonenumber" binding:"required"`
			Password    string `json:"password" binding:"required"`
			Role        string `json:"role" binding:"required"`
		}
	//Get Staff details
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill blanks"})
		return
	}
	//check if staff email alredy exist
	var existingstaff models.Admin
	if err := database.DB.Find("email=?", data.Email).First(&existingstaff).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email alredy exist"})
		return
	}
	//checking emaail
	if !utils.EmailCheck(data.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format "})
		return
	}
	//veryfy phone Number
	if !utils.PhoneNumberCheck(data.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone Number"})
		return
	}
	//password Secure
	hash, err := utils.GenerateHash(data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed while hashing password"})
		return
	}
	// staff details to store on database
	staff := models.Admin{
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    hash,
		RoleID:      constants.RoleStaff,
	}
	//Storing new staff data to database
	if err := database.DB.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff"})
		return
	}
	//sucess responce
	c.JSON(http.StatusOK, gin.H{"message": "Staff created sucessfully",
		"staff": staff})
}
