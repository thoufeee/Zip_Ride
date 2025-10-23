package prizepoolmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// update prize pool

func UpdatePrizePool(c *gin.Context) {
	id := c.Param("id")

	var prizepool models.PrizePool

	if err := database.DB.Where("id = ?", id).First(&prizepool).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "no prize pool found"})
		return
	}

	var input struct {
		VehicleType string   `json:"vehicle_type"`
		Commission  *float64 `json:"commission"`
		BonusAmount *float64 `json:"bonusamount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	if input.VehicleType != "" {
		prizepool.VehicleType = input.VehicleType
	}

	if input.Commission != nil {
		prizepool.Commission = *input.Commission
	}

	if input.BonusAmount != nil {
		prizepool.BonusAmount = *input.BonusAmount
	}

	if err := database.DB.Save(&prizepool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to update prize pool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "succsessfuly updated prize pool"})
}
