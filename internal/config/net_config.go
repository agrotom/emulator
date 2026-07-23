package config

import "time"

type NetConfig struct {
	MaxTimeout        time.Duration `json:"maxTimeout"`
	MaxReconnectTries int           `json:"maxReconnectTries"`
	MaxResendTries    int           `json:"maxResendTries"`
	ResendWaitTime    time.Duration `json:"resendWaitTime"`
}
