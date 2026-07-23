package track

import (
	"math"
	"testing"

	"github.com/agrotom/emulator/internal/mathutil"
)

var StartCoords = mathutil.Vector2f{X: 51.546529, Y: 46.037097}
var EndCoords = mathutil.Vector2f{X: 51.536257, Y: 46.02405}

func TestOSRMResponse(t *testing.T) {
	coords, err := GetOSRMRoute(StartCoords, EndCoords)

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	for _, coord := range coords {
		t.Logf("%f, %f", coord.X, coord.Y)
	}
}

func TestInterpolate(t *testing.T) {
	coords, err := GetOSRMRoute(StartCoords, EndCoords)

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	interpCoords := Interpolate(coords, StepM)

	for _, coord := range interpCoords {
		t.Logf("%f, %f", coord.X, coord.Y)
	}
}

func TestSimulating(t *testing.T) {
	coords, err := GetOSRMRoute(StartCoords, EndCoords)

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	interpCoords := Interpolate(coords, StepM)

	records := Simulate(interpCoords, 1)

	for _, rd := range records {
		t.Logf("%+v", rd)
	}
}

func TestOSRMNearestPoint(t *testing.T) {
	vec := mathutil.Vector2f{
		X: 55.49,
		Y: 36.80,
	}

	nearVec, err := GetOSRMNearestPoint(vec)

	if err != nil {
		t.Errorf("error while getting osrm nearest point: %s", err.Error())
		return
	}

	if nearVec.X != 55.490003 || nearVec.Y != 36.816152 {
		t.Errorf("near vec is wrong: original (%f, %f) and result (%f, %f)", vec.X, vec.Y, nearVec.X, nearVec.Y)
		return
	}
}

func TestDestination45Point(t *testing.T) {
	vec := mathutil.Vector2f{
		X: 55.49,
		Y: 36.80,
	}

	nearVec := Destination(vec, 2000, 45*math.Pi/180)

	if nearVec.X <= vec.X || nearVec.Y <= vec.Y {
		t.Errorf("near vec is wrong: original (%f, %f) and result (%f, %f)", vec.X, vec.Y, nearVec.X, nearVec.Y)
		return
	}
}

func TestDestination135Point(t *testing.T) {
	vec := mathutil.Vector2f{
		X: 55.49,
		Y: 36.80,
	}

	nearVec := Destination(vec, 2000, 135*math.Pi/180)

	if nearVec.X >= vec.X || nearVec.Y <= vec.Y {
		t.Errorf("near vec is wrong: original (%f, %f) and result (%f, %f)", vec.X, vec.Y, nearVec.X, nearVec.Y)
		return
	}
}

func TestDestination210Point(t *testing.T) {
	vec := mathutil.Vector2f{
		X: 55.49,
		Y: 36.80,
	}

	nearVec := Destination(vec, 2000, 210*math.Pi/180)

	if nearVec.X >= vec.X || nearVec.Y >= vec.Y {
		t.Errorf("near vec is wrong: original (%f, %f) and result (%f, %f)", vec.X, vec.Y, nearVec.X, nearVec.Y)
		return
	}
}

func TestDestination300Point(t *testing.T) {
	vec := mathutil.Vector2f{
		X: 55.49,
		Y: 36.80,
	}

	nearVec := Destination(vec, 2000, 300*math.Pi/180)

	if nearVec.X <= vec.X || nearVec.Y >= vec.Y {
		t.Errorf("near vec is wrong: original (%f, %f) and result (%f, %f)", vec.X, vec.Y, nearVec.X, nearVec.Y)
		return
	}
}
