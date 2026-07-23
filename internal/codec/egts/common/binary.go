package common

type Binary interface {
	Encode() ([]byte, error)
	Decode([]byte) error
	Length() uint16
}
