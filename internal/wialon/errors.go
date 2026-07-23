package wialon

import "errors"

var (
	ErrConnection   = errors.New("server connection error")
	ErrSetDeadline  = errors.New("setting deadline of ctx error")
	ErrDisconnected = errors.New("not connected to server")
	ErrTimeout      = errors.New("timeout")
	ErrNoReader     = errors.New("reader is nil")
)
