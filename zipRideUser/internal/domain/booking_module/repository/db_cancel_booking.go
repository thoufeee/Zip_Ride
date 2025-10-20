package repository

import (
	"fmt"
	"zipride/database"
	"zipride/internal/models"
)

func CancelBooking(bookingID uint, reason string) error {
	result := database.DB.Model(&models.Booking{}).
		Where("id = ? AND status NOT IN ?", bookingID, []string{models.StatusCancelled, models.StatusCompleted}).
		Updates(map[string]interface{}{
			"status":        models.StatusCancelled,
			"cancel_reason": reason,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking cannot be cancelled (already cancelled or completed)")
	}

	return nil
}
