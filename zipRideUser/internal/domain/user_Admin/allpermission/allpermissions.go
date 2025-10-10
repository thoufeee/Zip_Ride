package allpermission

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// show all permissions

func Permissions(c *gin.Context) {

	var permissions []models.Permission

	if err := database.DB.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no record found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": permissions})
}
