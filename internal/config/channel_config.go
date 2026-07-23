package config

import "time"

type ChannelConfig struct {
	Jitter time.Duration `json:"jitter"`
	Delay  time.Duration `json:"delay"`

	LossPercent      float32 `json:"lossPercent"`
	ConnBreakPercent float32 `json:"connBreakPercent"`

	MaxLossCount      int `json:"maxLossCount"`
	MaxConnBreakCount int `json:"maxConnBreakCount"`
}
