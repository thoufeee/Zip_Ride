package dashbord

import (
	"net/http"
	"time"
	"zipride/database"

	"github.com/gin-gonic/gin"
)
//struct for booking count
type BookingCount struct {
	Date  string `json:"date"`
	Count uint   `json:"count"`
}
//booking per dat count
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
	for i := range perday {
		t, _ := time.Parse("2006-01-02", perday[i].Date)
		perday[i].Date = t.Format("Jan 02")
	}
	//sucess responce and result
	c.JSON(http.StatusOK, gin.H{
		"message": "Daily booking counts fetched successfully",
		"data":    perday,
	})
}
