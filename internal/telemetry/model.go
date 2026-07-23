package telemetry

import "time"

type NavRecord struct {
	DeltaTime time.Time `json:"delta"`
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"lon"`
	Speed     int       `json:"speed"`
	Course    int       `json:"course"`
}
