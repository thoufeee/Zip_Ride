package prizepoolmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// delete prize pool
func DeletePrizePool(c *gin.Context) {
	id := c.Param("id")

	var prizepool models.PrizePool

	if err := database.DB.Where("id = ?", id).First(&prizepool).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "prize pool not found"})
		return
	}

	if err := database.DB.Delete(&prizepool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to delete prize pool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfully Deleted Prize Pool"})
}
