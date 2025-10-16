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
	var Staffs []models.Admin

	//fetch staff data from database
	if err := database.DB.Where("role_id = ?", constants.Staff).Find(&Staffs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to fetch staffs"})
		return
	}

	//to check the staffs length
	if len(Staffs) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No staffs found", "Data": []models.Admin{}})
		return
	}

	// sucess responce
	c.JSON(http.StatusOK, gin.H{"res": Staffs})
}
