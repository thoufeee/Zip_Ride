package services

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// ListDrivers lists drivers filtered by optional status
func ListDrivers(c *gin.Context) {
	status := c.Query("status")
	var drivers []models.Driver
	q := database.DB
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at desc").Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch drivers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"drivers": drivers})
}

// GetDriverDocs returns driver documents
func GetDriverDocs(c *gin.Context) {
	id := c.Param("id")
	var docs models.DriverDocuments
	if err := database.DB.Where("driver_id = ?", id).First(&docs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "docs not found"})
		return
	}
	c.JSON(http.StatusOK, docs)
}
