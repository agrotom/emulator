package wialon

import (
	"context"

	"github.com/agrotom/emulator/internal/codec"
)

type Client interface {
	Close() error
	Dial(ctx context.Context) error
	SendPacket(ctx context.Context, packet codec.Packet) ([]byte, error)
}
