package wialonips

import "errors"

var (
	ErrNoReader    = errors.New("reader is nil")
	ErrReading     = errors.New("reading error")
	ErrBadResponse = errors.New("bad response")
)
