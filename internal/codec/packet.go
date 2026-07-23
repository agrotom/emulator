package codec

type Packet interface {
	Encode() ([]byte, error)
	ValidateResponse([]byte) error
}