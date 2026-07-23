package telemetry

import "math/rand/v2"

type BoundingBox struct {
	MinLat float64 `json:"minLat"`
	MaxLat float64 `json:"maxLat"`
	MinLon float64 `json:"minLon"`
	MaxLon float64 `json:"maxLon"`
}

type BoundingBoxCollection []BoundingBox

func (col BoundingBoxCollection) ChooseRandomBounds() BoundingBox {
	return col[rand.IntN(len(col))]
}

var (
	MoscowBoundingBox BoundingBox = BoundingBox{
		MinLat: 55.49,
		MaxLat: 56.02,
		MinLon: 36.80,
		MaxLon: 38.10,
	}
)
