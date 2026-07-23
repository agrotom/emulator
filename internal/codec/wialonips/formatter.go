package wialonips

import (
	"fmt"
	"math"
	"time"

	"github.com/agrotom/emulator/internal/mathutil"
)

func WriteDate(fstr *string, date *time.Time) {
	if date == nil {
		*fstr += "NA;"
	} else {
		*fstr += date.Format("020106") + ";"
	}
}

func WriteTime(fstr *string, time *time.Time) {
	if time == nil {
		*fstr += "NA;"
	} else {
		*fstr += time.Format("150405") + ";"
	}
}

func FormatLatCoord(coord float64) (string, string) {
	var dir string

	if coord >= 0 {
		dir = "N"
	} else {
		dir = "S"
	}

	coord = math.Abs(coord)

	deg := math.Floor(coord)
	minutes := (coord - deg) * 60

	return fmt.Sprintf("%02d%07.4f", int(deg), minutes), dir
}

func FormatLonCoord(coord float64) (string, string) {
	var dir string

	if coord >= 0 {
		dir = "E"
	} else {
		dir = "W"
	}

	coord = math.Abs(coord)

	deg := math.Floor(coord)
	minutes := (coord - deg) * 60

	return fmt.Sprintf("%02d%07.4f", int(deg), minutes), dir
}

func WriteVec2f(fstr *string, vec *mathutil.Vector2f) {

	if vec == nil {
		*fstr += "NA;NA;"
	} else {
		*fstr += fmt.Sprintf("%g;%g", vec.X, vec.Y) + ";"
	}
}

func WriteInt(fstr *string, val *int) {
	if val == nil {
		*fstr += "NA;"
	} else {
		*fstr += fmt.Sprintf("%d", *val) + ";"
	}
}

func WriteFloat64(fstr *string, val *float64) {
	if val == nil {
		*fstr += "NA;"
	} else {
		*fstr += fmt.Sprintf("%g", *val) + ";"
	}
}

func WriteString(fstr *string, val *string) {
	if val == nil {
		*fstr += ";"
	} else {
		*fstr += *val + ";"
	}
}

func WriteStringNA(fstr *string, val *string) {
	if val == nil {
		*fstr += "NA;"
	} else {
		*fstr += *val + ";"
	}
}
