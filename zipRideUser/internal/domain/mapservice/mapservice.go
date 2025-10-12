package mapservice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// struct for responce the distance and the time duration
type RouteResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	} `json:"routes"`
}

func GetRouteDistance(pickupLat, pickupLong, dropLat, dropLong float64) (float64, float64, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=false", pickupLong, pickupLat, dropLong, dropLat)
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
	return route.Routes[0].Distance / 1000, route.Routes[0].Duration / 60, nil

}
