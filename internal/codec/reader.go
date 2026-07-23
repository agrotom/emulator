package codec

import "io"

type PacketReader interface {
	Read() ([]byte, error)
	SetReader(io.Reader)
}
