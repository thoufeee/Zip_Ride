package repository

import (
	"fmt"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
)

// cancell bookinf only when status is pending
func CancelBooking(bookingID uint, reason string) error {
	// Only update if status is pending
	result := database.DB.Model(&models.Booking{}).
		Where("id = ? AND status = ?", bookingID, constants.StatusPending).
		Updates(map[string]interface{}{
			"status":        constants.StatusCancelled,
			"cancel_reason": reason,
		})
		//when an error is occured
	if result.Error != nil {
		return result.Error
	}
	// error message explaining why the cancellation failed.
	if result.RowsAffected == 0 { //how many rows in the database were actually updated.
		return fmt.Errorf("booking cannot be cancelled (status is not pending)")
	}

	return nil
}
