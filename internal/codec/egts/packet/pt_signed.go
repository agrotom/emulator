package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const PtSignedOffset = 2

type PtSignedAppData struct {
	SignatureLength uint16
	SignatureData   []byte
	PtAppData
}

func (p *PtSignedAppData) Decode(data []byte) error {
	var (
		err          error
		uint16Buffer = make([]byte, 2)
	)

	buf := bytes.NewReader(data)

	if _, err = buf.Read(uint16Buffer); err != nil {
		return fmt.Errorf("error while decoding response packet id: %w", err)
	}

	p.SignatureLength = binary.LittleEndian.Uint16(uint16Buffer)

	p.SignatureData = make([]byte, p.SignatureLength)

	if _, err = buf.Read(p.SignatureData); err != nil {
		return fmt.Errorf("error while decoding response packet id: %w", err)
	}

	p.PtAppData = PtAppData{}
	if err = p.PtAppData.Decode(data[PtSignedOffset+p.SignatureLength:]); err != nil {
		return fmt.Errorf("error while decoding sdr: %w", err)
	}

	return nil
}

func (p *PtSignedAppData) Encode() ([]byte, error) {
	var (
		err      error
		sdrbytes []byte
	)

	buf := new(bytes.Buffer)

	if p.SignatureLength == 0 {
		p.SignatureLength = uint16(len(p.SignatureData))
	}

	if err = binary.Write(buf, binary.LittleEndian, p.SignatureLength); err != nil {
		return nil, fmt.Errorf("error while encoding signature length: %w", err)
	}

	if p.SignatureData != nil {
		if _, err = buf.Write(p.SignatureData); err != nil {
			return nil, fmt.Errorf("error while encoding signature data: %w", err)
		}
	}

	if p.PtAppData != nil {
		if sdrbytes, err = p.PtAppData.Encode(); err != nil {
			return nil, fmt.Errorf("error while encoding sdr: %w", err)
		}

		buf.Write(sdrbytes)
	}

	return buf.Bytes(), nil
}

func (p *PtSignedAppData) Length() uint16 {
	if ptBytes, err := p.Encode(); err != nil {
		return 0
	} else {
		return uint16(len(ptBytes))
	}
}
