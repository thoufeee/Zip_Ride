package allpermission

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// show all roles

func AllRoles(c *gin.Context) {

	var roles []models.Role

	if err := database.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no roles found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": roles})
}
