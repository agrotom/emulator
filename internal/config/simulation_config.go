package config

import (
	"time"

	"github.com/agrotom/emulator/internal/mathutil"
	"github.com/agrotom/emulator/internal/telemetry"
)

type SimulationConfig struct {
	Creds  Credentials `json:"creds"`
	NetCfg NetConfig   `json:"netCfg"`

	Protocol ProtocolType `json:"protocolType"`

	HasControlChannel bool          `json:"hasControlChannel"`
	HAProxy           bool          `json:"haproxy"`
	ChannelCfg        ChannelConfig `json:"channelCfg"`

	BoundingBoxes telemetry.BoundingBoxCollection `json:"boundingBoxes"`

	AutoStartPos bool              `json:"autoStartPos"`
	StartPos     mathutil.Vector2f `json:"startPos"`

	AutoEndPos bool              `json:"autoEndPos"`
	EndPos     mathutil.Vector2f `json:"endPos"`

	Sattelites int `json:"sats"`

	StepDistanceM int           `json:"stepDistanceM"`
	StepMillis    time.Duration `json:"stepMillis"`

	MinDist float64 `json:"minDist"`
	MaxDist float64 `json:"maxDist"`
}
