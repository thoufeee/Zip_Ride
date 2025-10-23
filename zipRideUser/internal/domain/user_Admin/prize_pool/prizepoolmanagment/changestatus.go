package prizepoolmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// change status

func UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var prizepool models.PrizePool

	if err := database.DB.Where("id = ?", id).First(&prizepool).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "prize pool not found"})
		return
	}

	var input struct {
		Active bool `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "inavlid request"})
		return
	}

	prizepool.Active = input.Active

	if err := database.DB.Save(&prizepool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly Updated Status"})
}
