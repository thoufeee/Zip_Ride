package controllers

import (
	"net/http"
	"strconv"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// unblock staff
func UnblockStaff(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid staff id"})
		return
	}

	var staff models.Admin

	if err := database.DB.Where("id = ? AND role = ?", id, constants.RoleAdmin).First(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff not found"})
		return
	}

	staff.Block = false

	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to save changes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "user unblocked successfuly"})
}
