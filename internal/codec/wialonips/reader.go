package wialonips

import (
	"bufio"
	"fmt"
	"io"
)

type WialonIPSReader struct {
	r *bufio.Reader
}

func CreateWialonReader() *WialonIPSReader {
	return &WialonIPSReader{}
}

func (wr *WialonIPSReader) Read() ([]byte, error) {
	if wr.r == nil {
		return nil, ErrNoReader
	}

	resp, err := wr.r.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReading, err)
	}

	return []byte(resp), nil
}

func (wr *WialonIPSReader) SetReader(r io.Reader) {
	wr.r = bufio.NewReader(r)
}
