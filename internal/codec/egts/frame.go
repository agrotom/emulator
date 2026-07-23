package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/codec/egts/packet"
)

const (
	DefaultHeaderLength byte = 11
	FullHeaderLength    byte = 16
)

type EGTSFrame struct {
	ProtocolVersion byte
	SecurityKeyID   byte

	Prefix     common.Prefix
	Route      bool
	Encryption common.Encryption
	Compressed bool
	Priority   common.Priority

	HeaderLength     byte
	HeaderEncoding   byte
	FrameDataLength  uint16
	PacketIdentifier uint16
	PacketType       common.PacketType

	PeerAddress      uint16
	RecipientAddress uint16
	TimeToLive       byte

	HeaderCheckSum byte

	ServicesFrameData         common.Binary
	ServicesFrameDataCheckSum uint16
}

func (frame *EGTSFrame) Decode(data []byte) error {
	var (
		err          error
		flags        byte
		uint16Buffer = make([]byte, 2)
	)

	buf := bytes.NewReader(data)

	if frame.ProtocolVersion, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding protocol version: %w", err)
	}

	if frame.SecurityKeyID, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding security key: %w", err)
	}

	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding flags: %w", err)
	}

	frame.Prefix = common.Prefix(common.Bits(flags, 6, 2))
	frame.Route = common.Bit(flags, 5)
	frame.Encryption = common.Encryption(common.Bits(flags, 3, 2))
	frame.Compressed = common.Bit(flags, 2)
	frame.Priority = common.Priority(common.Bits(flags, 0, 2))

	if frame.HeaderLength, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding header length: %w", err)
	}

	if frame.HeaderEncoding, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding header encoding: %w", err)
	}

	if _, err := buf.Read(uint16Buffer); err != nil {
		return fmt.Errorf("error while decoding frame data length: %w", err)
	}

	frame.FrameDataLength = binary.LittleEndian.Uint16(uint16Buffer)

	if _, err := buf.Read(uint16Buffer); err != nil {
		return fmt.Errorf("error while decoding packet identifier: %w", err)
	}

	frame.PacketIdentifier = binary.LittleEndian.Uint16(uint16Buffer)

	if pt, err := buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding packet type: %w", err)
	} else {
		frame.PacketType = common.PacketType(pt)
	}

	if frame.Route {
		if _, err := buf.Read(uint16Buffer); err != nil {
			return fmt.Errorf("error while decoding peer address: %w", err)
		}

		frame.PeerAddress = binary.LittleEndian.Uint16(uint16Buffer)

		if _, err := buf.Read(uint16Buffer); err != nil {
			return fmt.Errorf("error while decoding recipient address: %w", err)
		}

		frame.RecipientAddress = binary.LittleEndian.Uint16(uint16Buffer)

		if frame.TimeToLive, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("error while decoding TTL: %w", err)
		}
	}

	if frame.HeaderCheckSum, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding header check sum: %w", err)
	}

	if frame.HeaderCheckSum != common.CRC8(data[:frame.HeaderLength-1]) {
		return fmt.Errorf("error while decoding egts frame: header check sum doesn't coincide")
	}

	dataFrameBytes := make([]byte, frame.FrameDataLength)

	if _, err := buf.Read(dataFrameBytes); err != nil {
		return fmt.Errorf("error while decoding frame data: %w", err)
	}

	if frame.ServicesFrameData, err = packet.Factory(frame.PacketType); err != nil {
		return fmt.Errorf("error while decoding frame data: %w", err)
	}

	if err = frame.ServicesFrameData.Decode(dataFrameBytes); err != nil {
		return fmt.Errorf("error while decoding frame data: %w", err)
	}

	if _, err = buf.Read(uint16Buffer); err != nil {
		return fmt.Errorf("error while decoding service data frame checksum: %w", err)
	}

	frame.ServicesFrameDataCheckSum = binary.LittleEndian.Uint16(uint16Buffer)

	if frame.ServicesFrameDataCheckSum != common.CRC16(dataFrameBytes) {
		return fmt.Errorf("error while decoding egts frame: service data frame check sum doesn't coincide")
	}

	return nil
}

func (frame *EGTSFrame) Encode() ([]byte, error) {
	var (
		err   error
		flags byte
	)

	buf := new(bytes.Buffer)

	if err = buf.WriteByte(frame.ProtocolVersion); err != nil {
		return nil, fmt.Errorf("error while encoding protocol version: %w", err)
	}

	if err = buf.WriteByte(frame.SecurityKeyID); err != nil {
		return nil, fmt.Errorf("error while encoding security key: %w", err)
	}

	common.SetBits(&flags, 6, 2, byte(frame.Prefix))
	common.SetBit(&flags, 5, frame.Route)
	common.SetBits(&flags, 3, 2, byte(frame.Encryption))
	common.SetBit(&flags, 2, frame.Compressed)
	common.SetBits(&flags, 0, 2, byte(frame.Priority))

	if err = buf.WriteByte(flags); err != nil {
		return nil, fmt.Errorf("error while encoding flags: %w", err)
	}

	if frame.HeaderLength == 0 {
		if frame.Route {
			frame.HeaderLength = FullHeaderLength
		} else {
			frame.HeaderLength = DefaultHeaderLength
		}
	}

	if err = buf.WriteByte(frame.HeaderLength); err != nil {
		return nil, fmt.Errorf("error while encoding header length: %w", err)
	}

	if err = buf.WriteByte(frame.HeaderEncoding); err != nil {
		return nil, fmt.Errorf("error while encoding header encoding: %w", err)
	}

	sfdBytes, err := frame.ServicesFrameData.Encode()

	if err != nil {
		return nil, fmt.Errorf("error while encoding services frame data: %w", err)
	}

	frame.FrameDataLength = uint16(len(sfdBytes))

	if err = binary.Write(buf, binary.LittleEndian, frame.FrameDataLength); err != nil {
		return nil, fmt.Errorf("error while encoding packet id: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, frame.PacketIdentifier); err != nil {
		return nil, fmt.Errorf("error while encoding packet id: %w", err)
	}

	if err = buf.WriteByte(byte(frame.PacketType)); err != nil {
		return nil, fmt.Errorf("error while encoding packet type: %w", err)
	}

	if frame.Route {
		if err = binary.Write(buf, binary.LittleEndian, frame.PeerAddress); err != nil {
			return nil, fmt.Errorf("error while encoding peer address: %w", err)
		}

		if err = binary.Write(buf, binary.LittleEndian, frame.RecipientAddress); err != nil {
			return nil, fmt.Errorf("error while encoding recipient address: %w", err)
		}

		if err = buf.WriteByte(frame.TimeToLive); err != nil {
			return nil, fmt.Errorf("error while encoding TTL: %w", err)
		}
	}

	frame.HeaderCheckSum = common.CRC8(buf.Bytes())
	buf.WriteByte(frame.HeaderCheckSum)

	if _, err = buf.Write(sfdBytes); err != nil {
		return nil, fmt.Errorf("error while writing services frame data: %w", err)
	}

	frame.ServicesFrameDataCheckSum = common.CRC16(sfdBytes)

	if err = binary.Write(buf, binary.LittleEndian, frame.ServicesFrameDataCheckSum); err != nil {
		return nil, fmt.Errorf("error while encoding service data frame checksum: %w", err)
	}

	return buf.Bytes(), nil
}
