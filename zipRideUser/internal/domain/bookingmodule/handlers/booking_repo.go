package bookingmodule

import (
	"zipride/internal/models"

	"gorm.io/gorm"
)
//DB helper for booking model
type BookingRepo struct {
    DB *gorm.DB
}

//fetches a booking by ID
func NewBookingRepo(db *gorm.DB) *BookingRepo {
    return &BookingRepo{DB: db}
}
//saves a booking
func (r *BookingRepo) CreateBooking(b *models.Booking) error {
    return r.DB.Create(b).Error
}
//constructor to create the repo
func (r *BookingRepo) GetBookingByID(id uint) (*models.Booking, error) {
    var booking models.Booking
    err := r.DB.First(&booking, id).Error
    return &booking, err
}