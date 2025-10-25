package dashbord

import (
	"net/http"
	"time"
	"zipride/database"

	"github.com/gin-gonic/gin"
)

// MonthlyRevenue represents revenue per month
type MonthlyRevenue struct {
	MonthStart time.Time `json:"month_start"`
	Revenue    float64   `json:"revenue"`
}

// RevenuePerMonth fetches revenue for completed bookings per month
func RevenuePerMonth(c *gin.Context) {
	var revenues []MonthlyRevenue 		//to store the fetching data
	//query to get the monthly data
	query := `
		SELECT 
			DATE_TRUNC('month', created_at) AS month_start,
			COALESCE(SUM(fare),0) AS revenue
		FROM bookings
		WHERE status = 'completed'
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at);
	`
	//database integration based on query
	if err := database.DB.Raw(query).Scan(&revenues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch monthly revenue",
			"error":   err.Error(),
		})
		return
	}

	// Format month as YYYY-MM
	formatted := make([]map[string]interface{}, len(revenues))
	for i, r := range revenues {
		formatted[i] = map[string]interface{}{
			"month":   r.MonthStart.Format("2006-01"), // e.g., "2025-10"
			"revenue": r.Revenue,
		}
	}
	//sucess reponce
	c.JSON(http.StatusOK, gin.H{
		"message": "Monthly revenue for completed bookings fetched successfully",
		"data":    formatted,
	})
}
