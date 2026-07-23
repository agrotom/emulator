package channel

import "time"

type ControlChannelOptions struct {
	Jitter time.Duration
	Delay  time.Duration

	LossPercent      float32
	ConnBreakPercent float32

	MaxLossCount      int
	MaxConnBreakCount int
}

func DefaultOptions() ControlChannelOptions {
	return ControlChannelOptions{
		Jitter:            0,
		Delay:             0,
		LossPercent:       0.0,
		ConnBreakPercent:  0.0,
		MaxLossCount:      0,
		MaxConnBreakCount: 0,
	}
}
