package controllers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// Get All Staffss
func GETAllStaff(c *gin.Context) {
	var admin []models.Admin

	//fetch staff data from database
	if err := database.DB.Where("role = ?", constants.RoleAdmin).Find(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to fetch"})
		return
	}

	//to check the staffs length
	if len(admin) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No staffs found", "Data": []models.Admin{}})
		return
	}

	// sucess responce
	c.JSON(http.StatusOK, gin.H{"res": admin})
}
