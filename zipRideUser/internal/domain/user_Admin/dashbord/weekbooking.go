package dashbord

import (
	"net/http"
	"time"
	"zipride/database"

	"github.com/gin-gonic/gin"
)

// for week struct
type WeeklyBooking struct {
	WeekStart time.Time `json:"week_start"`
	Count     uint      `json:"count"`
}

// boooking per week
func BookingPerWeek(c *gin.Context) {
	var perWeek []WeeklyBooking
	//query for getting booking per week
	query := `
		SELECT 
			DATE_TRUNC('week', created_at) AS week_start,
			COUNT(*) AS count
		FROM bookings
		WHERE status != 'cancelled'
		GROUP BY DATE_TRUNC('week', created_at)
		ORDER BY DATE_TRUNC('week', created_at);
	`

	if err := database.DB.Raw(query).Scan(&perWeek).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch weekly bookings",
			"error":   err.Error(),
		})
		return
	}

	// Return plain week_start as YYYY-MM-DD
	formatted := make([]map[string]interface{}, len(perWeek))
	for i, b := range perWeek {
		formatted[i] = map[string]interface{}{
			"week_start": b.WeekStart.Format("2006-01-02"),
			"count":      b.Count,
		}
	}
	//sucess responce
	c.JSON(http.StatusOK, gin.H{
		"message": "Weekly booking counts fetched successfully",
		"data":    formatted,
	})
}
