package prizepoolmanagment

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// for creating prize pools

func CreatePrizePool(c *gin.Context) {

	var input struct {
		VehicleType string  `json:"vehicle_type" binding:"required"`
		Commission  float64 `json:"commission" binding:"required"`
		BonusAmount float64 `json:"bonusamount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	if input.BonusAmount == 0 {
		input.BonusAmount = 0
	}

	newprize_pool := models.PrizePool{
		ID:          uuid.New(),
		VehicleType: input.VehicleType,
		Commission:  input.Commission,
		BonusAmount: input.BonusAmount,
		Active:      true,
	}

	if err := database.DB.Create(&newprize_pool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to create new prize pool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Successfuly Created New Prize Pool"})
}
