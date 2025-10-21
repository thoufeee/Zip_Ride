package vehiclemanagement

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// Fetch all vehicle fares
func GetAllVehicleFares(c *gin.Context) {
	var fares []models.Vehicle
	if err := database.DB.Find(&fares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fares"})
		return
	}

	c.JSON(http.StatusOK, fares)
}

// VehicleFareCreation
func VehicleFareCreation(c *gin.Context) {
	var fare models.Vehicle

	// Bind JSON body
	if err := c.ShouldBindJSON(&fare); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation
	if fare.VehicleType == "" || fare.BaseFare <= 0 || fare.PerKmRate <= 0 || fare.PeopleCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fare fields are required and must be valid"})
		return
	}

	// Save to database
	if err := database.DB.Create(&fare).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vehicle fare"})
		return
	}
	//sucess responce
	c.JSON(http.StatusCreated, gin.H{
		"message": "Vehicle fare added successfully",
		"data":    fare,
	})
}

// Get a single fare by ID
func GetVehicleFareByID(c *gin.Context) {
	id := c.Param("id")
	var fare models.Vehicle

	if err := database.DB.First(&fare, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle fare not found"})
		return
	}
	c.JSON(http.StatusOK, fare)
}

// Update a vehicle fare
func UpdateVehicleFare(c *gin.Context) {
	id := c.Param("id")
	var fare models.Vehicle

	if err := database.DB.First(&fare, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle fare not found"})
		return
	}

	var input models.Vehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Partial updates are allowed
	if err := database.DB.Model(&fare).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle fare"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vehicle fare updated successfully",
		"data":    fare,
	})
}

// Delete a fare
func DeleteVehicleFare(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Vehicle{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle fare"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehicle fare deleted successfully"})
}
