package configuration

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// all configurations

func AllConfiguration(c *gin.Context) {
	var config models.WebConfig

	if err := database.DB.Find(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "configuration not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": config})
}
