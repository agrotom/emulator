package track

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/agrotom/emulator/internal/mathutil"
)

const OSRMUrl = "https://osrm.test1.opvk.tech"

func GetOSRMNearestPoint(vec mathutil.Vector2f) (mathutil.Vector2f, error) {
	url := fmt.Sprintf("%s/nearest/v1/driving/%f,%f", OSRMUrl, vec.Y, vec.X)
	resp, err := http.Get(url)

	if err != nil {
		return vec, fmt.Errorf("%w: %w", ErrOSRMRequest, err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return vec, fmt.Errorf("%w: %w", ErrOSRMReading, err)
	}

	data := OSRMNearestResponse{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		return vec, err
	}

	vec.X = data.Waypoints[0].Location[1]
	vec.Y = data.Waypoints[0].Location[0]

	return vec, nil
}

func GetOSRMRoute(start mathutil.Vector2f, end mathutil.Vector2f) ([]mathutil.Vector2f, error) {
	coords := fmt.Sprintf("%f,%f;%f,%f", start.Y, start.X, end.Y, end.X)
	url := fmt.Sprintf("%s/route/v1/driving/%s?overview=full&geometries=geojson", OSRMUrl, coords)

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOSRMRequest, err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOSRMReading, err)
	}

	data := OSRMResponse{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	result := make([]mathutil.Vector2f, 0)

	for _, coords := range data.Routes[0].Geometry.Coordinates {
		result = append(result, mathutil.Vector2f{X: coords[1], Y: coords[0]})
	}

	return result, nil
}
