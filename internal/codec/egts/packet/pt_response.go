package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const PtResponseByteOffset = 3

type PtResponse struct {
	ResponsePacketID uint16
	ProcessingResult byte
	PtAppData
}

func (p *PtResponse) Decode(data []byte) error {
	var (
		err          error
		uint16Buffer = make([]byte, 2)
	)

	buf := bytes.NewReader(data)

	if _, err = buf.Read(uint16Buffer); err != nil {
		return fmt.Errorf("error while decoding response packet id: %w", err)
	}

	p.ResponsePacketID = binary.LittleEndian.Uint16(uint16Buffer)

	if p.ProcessingResult, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding processing result: %w", err)
	}

	p.PtAppData = PtAppData{}
	if err = p.PtAppData.Decode(data[PtResponseByteOffset:]); err != nil {
		return fmt.Errorf("error while decoding sdr: %w", err)
	}

	return nil
}

func (p *PtResponse) Encode() ([]byte, error) {
	var (
		err      error
		sdrbytes []byte
	)

	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, p.ResponsePacketID); err != nil {
		return nil, fmt.Errorf("error while encoding response packet id: %w", err)
	}

	if err = buf.WriteByte(p.ProcessingResult); err != nil {
		return nil, fmt.Errorf("error while encoding processing result: %w", err)
	}

	if p.PtAppData != nil {
		if sdrbytes, err = p.PtAppData.Encode(); err != nil {
			return nil, fmt.Errorf("error while reading sdr: %w", err)
		}

		if _, err = buf.Write(sdrbytes); err != nil {
			return nil, fmt.Errorf("error while encoding sdr: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func (p *PtResponse) Length() uint16 {
	if ptBytes, err := p.Encode(); err != nil {
		return 0
	} else {
		return uint16(len(ptBytes))
	}
}
