package bookingmodule

import "strings"

// vehicle details
type Vehicle struct {
	ID         uint    `json:"id"`
	Type       string  `json:"type"`      // Car, Auto, Bike
	BaseFare   float64 `json:"base_fare"` // starting fare
	PerKmRate  float64 `json:"per_km"`
	PerMinRate float64 `json:"per_min"`
	MinFare    float64 `json:"min_fare"`
}

// vehicle fare and details
func GetAvalableVehicles() []Vehicle {
	return []Vehicle{
		{ID: 1, Type: "Bike", BaseFare: 10, PerKmRate: 5, PerMinRate: 1, MinFare: 20},
		{ID: 2, Type: "Auto", BaseFare: 15, PerKmRate: 8, PerMinRate: 1.5, MinFare: 30},
		{ID: 3, Type: "Car", BaseFare: 25, PerKmRate: 12, PerMinRate: 2, MinFare: 50},
	}
}


func GetVehicleByType(vehicleType string) Vehicle {
	for _, v := range GetAvalableVehicles() {
		if strings.EqualFold(v.Type, vehicleType) {
			return v
		}
	}
	// default to Car if type not found
	return Vehicle{Type: "Car", BaseFare: 25, PerKmRate: 12, PerMinRate: 2, MinFare: 50}
}
