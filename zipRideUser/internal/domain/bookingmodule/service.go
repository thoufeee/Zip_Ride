package bookingmodule

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
	"zipride/database"
	"zipride/internal/domain/mapservice"
	"zipride/internal/models"
)

type EstimateResponse struct {
	Type     string  `json:"type"`
	Fare     float64 `json:"fare"`
	ETA      string  `json:"eta"`
	Distance string  `json:"distance"`
	Duration string  `json:"duration"`
}

// helper: round to given decimal places
func round(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(val*pow) / pow
}

// get extimate fare for seperate vehicles
func GetFareEstimates(pickupLat, pickupLong, dropLat, dropLong float64) ([]EstimateResponse, error) {
	distance, duration, err := mapservice.GetRouteDistance(pickupLat, pickupLong, dropLat, dropLong)
	if err != nil {
		return nil, err
	}

	var estimates []EstimateResponse
	for _, v := range GetAvalableVehicles() {
		fare := v.BaseFare + (v.PerKmRate * distance) + (v.PerMinRate * duration)
		if fare < v.MinFare {
			fare = v.MinFare
		}

		estimates = append(estimates, EstimateResponse{
			Type:     v.Type,
			Fare:     round(fare, 2),
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
	v := GetVehicleByType(vehicleType)

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
	fare = round(fare, 2)
	// To generate otp while book a service
	otp, err := GenerateOTP()
	if err != nil {
		return nil, err
	}

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

// OTP Generate to while book a servcie
func GenerateOTP() (string, error) {
	b := make([]byte, 2) // 2 bytes = 0-65535
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d", int(b[0])<<8|int(b[1])%10000), nil
}
