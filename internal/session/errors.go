package session

import "errors"

var (
	ErrSendPacket      = errors.New("sending packet error")
	ErrNoClient        = errors.New("client is nil")
	ErrNoPacketBuilder = errors.New("packet builder is nil")
)
