package session

import (
	"time"

	"github.com/agrotom/emulator/internal/mathutil"
	"github.com/agrotom/emulator/internal/telemetry"
)

type SimulationOptions struct {
	StartPos *mathutil.Vector2f
	EndPos   *mathutil.Vector2f

	MaxPacketToSend int

	Sattelites int

	StepDistance int
	StepMillis   time.Duration

	MinDist float64
	MaxDist float64
	Bounds  telemetry.BoundingBoxCollection
}
