package controllers

import (
	"net/http"
	"strconv"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// block staff

func BlockStaff(c *gin.Context) {
	idstr := c.Param("id")

	if idstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff id required"})
		return
	}

	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid staff id"})
		return
	}

	var staff models.Admin

	if err := database.DB.Where("id = ? AND role_id = ?", uint(id), constants.Staff).First(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff not found"})
		return
	}

	staff.Block = true

	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to save changes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "user blocked successfuly"})
}
