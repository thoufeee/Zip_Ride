package prizepoolmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// all price pools
func GetAllPrizePool(c *gin.Context) {
	var prizepool []models.PrizePool

	if err := database.DB.Find(&prizepool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no data found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": prizepool})
}
