package codec

import (
	"github.com/agrotom/emulator/internal/telemetry"
)

type LoginOptions struct {
	IMEI     string
	Password string
}

type DataOptions struct {
	telemetry.NavRecord
	Sattelites int
}

type PacketBuilder interface {
	LoginPacket(*LoginOptions) Packet
	DataPacket(*DataOptions) Packet
}
