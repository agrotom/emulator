package track

import (
	"math"
	"math/rand"
	"time"

	"github.com/agrotom/emulator/internal/mathutil"
	"github.com/agrotom/emulator/internal/telemetry"
)

const StepM int = 100
const StepMS int64 = 1000

const EarthRadius float64 = 6371000

func RandomPoint(box telemetry.BoundingBox) mathutil.Vector2f {
	lat := box.MinLat + rand.Float64()*(box.MaxLat-box.MinLat)
	lon := box.MinLon + rand.Float64()*(box.MaxLon-box.MinLon)

	return mathutil.Vector2f{
		X: lat,
		Y: lon,
	}
}

func Destination(vec mathutil.Vector2f, dist, bearing float64) mathutil.Vector2f {
	lat1 := vec.X * math.Pi / 180
	lon1 := vec.Y * math.Pi / 180

	lat2 := math.Asin(
		math.Sin(lat1)*math.Cos(dist/EarthRadius) +
			math.Cos(lat1)*math.Sin(dist/EarthRadius)*math.Cos(bearing),
	)

	lon2 := lon1 + math.Atan2(
		math.Sin(bearing)*math.Sin(dist/EarthRadius)*math.Cos(lat1),
		math.Cos(dist/EarthRadius)-math.Sin(lat1)*math.Sin(lat2),
	)

	return mathutil.Vector2f{
		X: lat2 * 180 / math.Pi,
		Y: lon2 * 180 / math.Pi,
	}
}

func RandomDestination(vec mathutil.Vector2f, minDist, maxDist float64) mathutil.Vector2f {
	if minDist > maxDist {
		temp := minDist
		minDist = maxDist
		maxDist = temp
	}

	dist := minDist + rand.Float64()*(maxDist-minDist)
	bearing := rand.Float64() * 2 * math.Pi

	lat1 := vec.X * math.Pi / 180
	lon1 := vec.Y * math.Pi / 180

	lat2 := math.Asin(
		math.Sin(lat1)*math.Cos(dist/EarthRadius) +
			math.Cos(lat1)*math.Sin(dist/EarthRadius)*math.Cos(bearing),
	)

	lon2 := lon1 + math.Atan2(
		math.Sin(bearing)*math.Sin(dist/EarthRadius)*math.Cos(lat1),
		math.Cos(dist/EarthRadius)-math.Sin(lat1)*math.Sin(lat2),
	)

	return mathutil.Vector2f{
		X: lat2 * 180 / math.Pi,
		Y: lon2 * 180 / math.Pi,
	}
}

func CalculateCourse(vec1 mathutil.Vector2f, vec2 mathutil.Vector2f) float64 {
	lat1, lon1 := vec1.X*math.Pi/180, vec1.Y*math.Pi/180
	lat2, lon2 := vec2.X*math.Pi/180, vec2.Y*math.Pi/180

	dlon := lon2 - lon1

	x := math.Sin(dlon) * math.Cos(lat2)
	y := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)

	return math.Mod(math.Atan2(x, y)*180.0/math.Pi+360.0, 360.0)
}

// Calculates in meters
func HaversineDistanceM(vec1 mathutil.Vector2f, vec2 mathutil.Vector2f) float64 {
	lat1, lon1 := vec1.X*math.Pi/180, vec1.Y*math.Pi/180
	lat2, lon2 := vec2.X*math.Pi/180, vec2.Y*math.Pi/180

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Pow(math.Sin(dlat/2), 2)
	a += math.Cos(lat1) * math.Cos(lat2) * math.Pow(math.Sin(dlon/2), 2)

	// math.Atan2 is used instead of math.Asin because it's more stable
	c := 2 * EarthRadius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c
}

func Interpolate(coords []mathutil.Vector2f, step int) []mathutil.Vector2f {
	result := make([]mathutil.Vector2f, 0)
	for i := 1; i < len(coords); i++ {
		pos1, pos2 := coords[i-1], coords[i]
		length := HaversineDistanceM(pos1, pos2)
		steps := max(2, int(length/float64(step)))

		for j := range steps {
			if i > 1 && j == 0 {
				continue
			}
			lat := pos1.X + (pos2.X-pos1.X)*float64(j)/(float64(steps)-1)
			lon := pos1.Y + (pos2.Y-pos1.Y)*float64(j)/(float64(steps)-1)
			result = append(result, mathutil.Vector2f{X: lat, Y: lon})
		}
	}

	return result
}

func Simulate(coords []mathutil.Vector2f, step time.Duration) []telemetry.NavRecord {
	records := []telemetry.NavRecord{}
	dt := time.Now().UTC()
	for i := 1; i < len(coords); i++ {
		pos1, pos2 := coords[i-1], coords[i]
		speed := 0
		if (i == len(coords)-1) || i == 1 {
			speed = 0
		} else {
			if rand.Float64() < 1.0/12.0 {
				speed = 0
			} else {
				speed = rand.Intn(111-60) + 60
			}
		}
		course := CalculateCourse(pos1, pos2)
		records = append(records, telemetry.NavRecord{DeltaTime: dt, Latitude: pos2.X, Longitude: pos2.Y, Speed: speed, Course: int(course)})
		dt = dt.Add(time.Millisecond * step)
	}

	return records
}
