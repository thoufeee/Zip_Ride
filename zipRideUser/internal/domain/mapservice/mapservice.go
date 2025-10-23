package mapservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// -------------------- Geocode --------------------
// Converts location name → latitude/longitude using OpenStreetMap
type GeoResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func Geocode(location string) (float64, float64, error) {
	endpoint := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", location)
	params.Add("format", "json")

	resp, err := http.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to call geocode API: %v", err)
	}
	defer resp.Body.Close()

	var results []GeoResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, fmt.Errorf("failed to parse geocode response: %v", err)
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("no results found for %s", location)
	}

	var lat, lon float64
	fmt.Sscanf(results[0].Lat, "%f", &lat)
	fmt.Sscanf(results[0].Lon, "%f", &lon)

	return lat, lon, nil
}

// -------------------- Route Distance --------------------
// Struct for OSRM response
type RouteResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	} `json:"routes"`
}

// GetRouteDistance calculates distance (km) & duration (minutes) between two coordinates
func GetRouteDistance(pickupLat, pickupLong, dropLat, dropLong float64) (float64, float64, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false",
		pickupLong, pickupLat, dropLong, dropLat,
	)

	res, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer res.Body.Close()

	var route RouteResponse
	if err := json.NewDecoder(res.Body).Decode(&route); err != nil {
		return 0, 0, err
	}

	if len(route.Routes) == 0 {
		return 0, 0, fmt.Errorf("no routes found")
	}

	distanceKm := route.Routes[0].Distance / 1000
	durationMin := route.Routes[0].Duration / 60
	return distanceKm, durationMin, nil
}
