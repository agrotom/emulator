package track

import "errors"

var (
	ErrOSRMRequest = errors.New("osrm request error")
	ErrOSRMReading = errors.New("osrm reading error")
)