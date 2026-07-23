package egts

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type FrameReader struct {
	r io.Reader
}

func NewReader(r io.Reader) *FrameReader {
	return &FrameReader{r}
}

func (fr *FrameReader) ReadFrame() (*EGTSFrame, error) {
	var err error

	header := make([]byte, DefaultHeaderLength)

	if _, err = io.ReadFull(fr.r, header); err != nil {
		return nil, fmt.Errorf("error while reading response packet from wialon: %w", err)
	}

	log.Println("Read some bytes")

	headerLength := header[3]

	if headerLength == FullHeaderLength {
		extra := make([]byte, FullHeaderLength-DefaultHeaderLength)

		if _, err = io.ReadFull(fr.r, extra); err != nil {
			return nil, fmt.Errorf("error while reading extra header info of response packet from wialon: %w", err)
		}

		header = append(header, extra...)
	} else if headerLength != DefaultHeaderLength {
		return nil, fmt.Errorf("invalid header length of response packet from wialon: %d", headerLength)
	}

	frameDataLength := binary.LittleEndian.Uint16(header[5:7])
	frameData := make([]byte, frameDataLength+2)

	if _, err = io.ReadFull(fr.r, frameData); err != nil {
		return nil, fmt.Errorf("error while reading frame data of response packet from wialon: %w", err)
	}

	packetBytes := append(header, frameData...)
	packet := &EGTSFrame{}

	if err = packet.Decode(packetBytes); err != nil {
		return nil, fmt.Errorf("error while decoding response packet from wialon: %w", err)
	}

	return packet, nil
}
