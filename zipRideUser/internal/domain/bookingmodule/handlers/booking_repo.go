package bookingmodule

import (
	"fmt"
	"zipride/internal/models"

	"gorm.io/gorm"
)

type BookingRepo struct {
	DB *gorm.DB
}

func NewBookingRepo(db *gorm.DB) *BookingRepo {
	return &BookingRepo{DB: db}
}

// Create a new booking
func (r *BookingRepo) CreateBooking(b *models.Booking) error {
	return r.DB.Create(b).Error
}

// Get booking by ID
func (r *BookingRepo) GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.DB.First(&booking, id).Error
	return &booking, err
}

// Get rides by user ID (history)
func (r *BookingRepo) GetRidesByUserID(userID uint) ([]models.Booking, error) {
	var rides []models.Booking
	err := r.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&rides).Error
	return rides, err
}

//  Cancel booking
func (r *BookingRepo) CancelBooking(bookingID, userID uint) error {
	// Only cancel if booking belongs to user and status is 'ongoing'
	result := r.DB.Model(&models.Booking{}).
		Where("id = ? AND user_id = ? AND status = ?", bookingID, userID, "pending").
		Update("status", "cancelled")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// No booking was updated → either invalid ID or not cancellable
		return fmt.Errorf("booking cannot be cancelled")
	}

	return nil
}
