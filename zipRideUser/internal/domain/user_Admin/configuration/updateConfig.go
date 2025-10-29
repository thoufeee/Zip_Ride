package configuration

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// update configuration

func UpdateConfig(c *gin.Context) {
	var config map[string]interface{}

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid data"})
		return
	}

	if len(config) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no fileds to update"})
		return
	}

	if err := database.DB.Model(&models.WebConfig{}).Where("id = ?", 1).Updates(config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Configuration updated Successfuly"})
}
