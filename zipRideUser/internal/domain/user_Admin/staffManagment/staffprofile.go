package staffmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/middleware"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// staff profile || manager || super admin
func StaffProfile(c *gin.Context) {

	staffID := middleware.GetUserID(c)

	if staffID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unathorized"})
		return
	}

	var staff models.Admin

	if err := database.DB.Find(&staff, staffID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": staff})
}
