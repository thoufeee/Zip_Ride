package dashbord

import (
	"net/http"
	"time"
	"zipride/database"

	"github.com/gin-gonic/gin"
)

// MonthlyBooking struct
type MonthlyBooking struct {
	MonthStart time.Time `json:"month_start"`  //store first day of month
	Count      uint      `json:"count"` //number count of booking
}

// BookingPerMonth returns booking count per month
func BookingPerMonth(c *gin.Context) {
	var perMonth []MonthlyBooking //slice to hold data
	//query month
	query := `
		SELECT 
			DATE_TRUNC('month', created_at) AS month_start,
			COUNT(*) AS count  
		FROM bookings
		WHERE status != 'cancelled'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at);
	`
	//database check and query check for month
	if err := database.DB.Raw(query).Scan(&perMonth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch monthly bookings",
			"error":   err.Error(),
		})
		return
	}

	// Format month_start as YYYY-MM-DD (first day of month)
	formatted := make([]map[string]interface{}, len(perMonth))
	for i, b := range perMonth {
		formatted[i] = map[string]interface{}{
			"month_start": b.MonthStart.Format("2006-01-02"),
			"count":       b.Count,
		}
	}
	//sucess responce
	c.JSON(http.StatusOK, gin.H{
		"message": "Monthly booking counts fetched successfully",
		"data":    formatted,
	})
}
