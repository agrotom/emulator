package wialonips

import (
	"github.com/agrotom/emulator/internal/codec"
)

type WialonPacketBuilder struct {
}

func CreateWialonPacketBuilder() *WialonPacketBuilder {
	return &WialonPacketBuilder{}
}

func (wpb *WialonPacketBuilder) LoginPacket(opts *codec.LoginOptions) codec.Packet {
	packet := &LoginPacket{
		IMEI:     opts.IMEI,
		Password: opts.Password,
	}

	return packet
}

func (wpb *WialonPacketBuilder) DataPacket(opts *codec.DataOptions) codec.Packet {
	packet := &DataPacket{
		Date:      &opts.DeltaTime,
		Time:      &opts.DeltaTime,
		Latitude:  &opts.Latitude,
		Longitude: &opts.Longitude,
		Speed:     &opts.Speed,
		Course:    &opts.Course,
		Sats:      &opts.Sattelites,
	}

	return packet
}
