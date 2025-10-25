package dashbord

import (
	"net/http"
	"time"
	"zipride/database"

	"github.com/gin-gonic/gin"
)

// struct for booking count
type BookingCount struct {
	Date  time.Time `json:"date"`
	Count uint   `json:"count"`
}

// booking per dat count
func BookingPerDay(c *gin.Context) {
	var perday []BookingCount
	//query for booking per day
	query := `
		SELECT 
			DATE(created_at) AS date,
			COUNT(*) AS count
		FROM bookings
		WHERE status != 'cancelled'  -- optional: exclude cancelled rides
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at);
	`
	//getting data from the database based on the query
	if err := database.DB.Raw(query).Scan(&perday).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch bookings per day",
			"error":   err.Error(),
		})
		return
	}

	// Format date to "Jan 02"
	formatted := make([]map[string]interface{}, len(perday))
	for i, b := range perday {
		formatted[i] = map[string]interface{}{
			"date":  b.Date.Format("Jan 02"), // formatted as "Oct 25"
			"count": b.Count,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daily booking counts fetched successfully",
		"data":    formatted,
	})
}
