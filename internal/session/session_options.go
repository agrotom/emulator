package session

import (
	"time"
)

type SessionOptions struct {
	UnitID   string
	Password string

	// Inflicts to reading and writing timeouts of connection in seconds.
	//
	// Provide -1 if there is an infinity timeout session (not recommended).
	MaxTimeout time.Duration

	MaxReconnectTries int

	MaxResendTries int

	ResendWaitTime time.Duration
}

// Creates default options with MaxTimeout=5 and MaxReconnectTries=5
func DefaultOptions(unitID, password string) SessionOptions {
	return SessionOptions{
		UnitID:            unitID,
		Password:          password,
		MaxTimeout:        5,
		MaxResendTries:    5,
		MaxReconnectTries: 5,
		ResendWaitTime:    5,
	}
}
