package track

type Geometry struct {
	Coordinates [][]float64 `json:"coordinates"`
}

type Route struct {
	Geometry Geometry `json:"geometry"`
}

type OSRMResponse struct {
	Code   string  `json:"code"`
	Routes []Route `json:"routes"`
}

type Waypoint struct {
	Location []float64 `json:"location"`
}

type OSRMNearestResponse struct {
	Code      string  `json:"code"`
	Waypoints []Waypoint `json:"waypoints"`
}
