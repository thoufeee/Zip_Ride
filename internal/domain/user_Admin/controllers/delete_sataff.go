package controllers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

func DeleteStaff(c *gin.Context) {
	//getting staff id from param
	staffid := c.Param("id")
	if staffid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "staff id Required"})
		return
	}
	//fetch staff from database
	var staff models.Admin
	if err := database.DB.First("id = ? AND RoleID = ?", staffid, constants.RoleStaff).Find(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "staff not found"})
		return
	}
	//delete staff data
	if err := database.DB.Delete(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff"})
		return
	}
	// sucess messge
	c.JSON(http.StatusOK, gin.H{"message": "Staff deleted sucessfully",})
}
