package bookingmodule

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"zipride/database"
	"zipride/internal/domain/bookingmodule"
	"zipride/internal/domain/mapservice"
	"zipride/internal/models"
	"zipride/utils"
)

type EstimateResponse struct {
	Type     string  `json:"type"`
	Fare     float64 `json:"fare"`
	ETA      string  `json:"eta"`
	Distance string  `json:"distance"`
	Duration string  `json:"duration"`
}

// get extimate fare for seperate vehicles
func GetFareEstimates(pickupLat, pickupLong, dropLat, dropLong float64) ([]EstimateResponse, error) {
	distance, duration, err := mapservice.GetRouteDistance(pickupLat, pickupLong, dropLat, dropLong)
	if err != nil {
		return nil, err
	}

	var estimates []EstimateResponse
	for _, v := range bookingmodule.GetAvalableVehicles() {
		fare := v.BaseFare + (v.PerKmRate * distance) + (v.PerMinRate * duration)
		if fare < v.MinFare {
			fare = v.MinFare
		}

		estimates = append(estimates, EstimateResponse{
			Type:     v.Type,
			Fare:     utils.Round(fare),
			ETA:      fmt.Sprintf("%.0f mins", duration/5), // dummy ETA
			Distance: fmt.Sprintf("%.2f km", distance),     // distance in km
			Duration: fmt.Sprintf("%.1f mins", duration),   // total trip time
		})
	}
	return estimates, nil
}

// create booking function
func CreateBooking(userID uint, pickupLat, pickupLong, dropLat, dropLong float64, vehicleType string, scheduleAt *time.Time) (*models.Booking, error) {

	// get vehicle struct
	v := bookingmodule.GetVehicleByType(vehicleType)

	// calculate distance & duration
	distance, duration, err := mapservice.GetRouteDistance(pickupLat, pickupLong, dropLat, dropLong)
	if err != nil {
		return nil, err
	}

	// calculate fare
	fare := v.BaseFare + (v.PerKmRate * distance) + (v.PerMinRate * duration)
	if fare < v.MinFare {
		fare = v.MinFare
	}
	fare = utils.Round(fare)
	// To generate otp while book a service
	otp := utils.GeneratorOtp()
	// create booking struct
	booking := &models.Booking{
		UserID:     userID,
		PickupLat:  pickupLat,
		PickupLong: pickupLong,
		DropLat:    dropLat,
		DropLong:   dropLong,
		Vehicle:    vehicleType,
		Fare:       fare,
		OTP:        otp,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if scheduleAt != nil {
		// Later booking -> store in Postgres
		booking.ScheduleAt = *scheduleAt
		if err := database.DB.Create(booking).Error; err != nil {
			return nil, err
		}
	} else {
		// Now booking -> store in Redis
		key := fmt.Sprintf("booking:%d", booking.ID)
		data, _ := json.Marshal(booking)
		err := database.RDB.Set(context.Background(), key, data, 30*time.Minute).Err()
		if err != nil {
			return nil, err
		}
	}

	return booking, nil
}

// helper function for date and time
func ParseUserSchedule(dateStr, timeStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	timeStr = strings.TrimSpace(strings.ToUpper(timeStr))

	combined := fmt.Sprintf("%s %s", dateStr, timeStr)
	layout := "2006-01-02 03:04 PM" // 12-hour clock with AM/PM

	t, err := time.Parse(layout, combined)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse combined date+time: %v", err)
	}
	return t, nil

}
