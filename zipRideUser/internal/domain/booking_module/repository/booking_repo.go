package repository

import (
	"time"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
)

func SaveBooking(b *models.Booking) error {
	return database.DB.Create(b).Error
}

// get booking by id
func GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := database.DB.First(&booking, id).Error
	return &booking, err
}

// repository/booking.go
func IsDuplicateBooking(userID uint, pickupLat, pickupLong, dropLat, dropLong float64, vehicle string, scheduleAt *time.Time) (bool, error) {
	var count int64
	db := database.DB

	query := db.Model(&models.Booking{}).
		Where("user_id = ? AND pickup_lat = ? AND pickup_long = ? AND drop_lat = ? AND drop_long = ? AND vehicle = ? AND status = ?",
			userID, pickupLat, pickupLong, dropLat, dropLong, vehicle, constants.StatusPending)

	if scheduleAt != nil {
		query = query.Where("schedule_at = ?", scheduleAt)
	} else {
		// for instant bookings, check last 2 minutes
		twoMinutesAgo := time.Now().Add(-2 * time.Minute)
		query = query.Where("created_at >= ?", twoMinutesAgo)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

//get users history function
func GetUserBookingHistory(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := database.DB.
		Model(&models.Booking{}).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&bookings).Error
	return bookings, err
}
