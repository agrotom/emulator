package teledata

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
)

const (
	SpeedMask uint16 = 0x3FFF
	ALTSMask  uint16 = 0x4000
	DIRHMask  uint16 = 0x8000
	DirMask   byte   = 0b10000000
)

type SrPosData struct {
	NavigationTime      uint32
	Latitude            uint32
	Longitude           uint32
	AltitudeFieldExists bool
	LongitudeHemisphere bool
	LatitudeHemisphere  bool
	Moving              bool
	BlackBox            bool
	CoordinateSystem    bool
	Fixture             bool
	Validness           bool
	Speed               uint16
	AltitudeSign        bool
	DirectionHighestBit bool
	Direction           byte
	Odometer            uint32
	DigitalInputs       byte
	Source              byte
	Altitude            uint32
	SourceData          uint16
}

func (sr *SrPosData) Decode(data []byte) error {
	var (
		flags byte
		err   error
	)

	buf := bytes.NewReader(data)

	if value, err := common.ReadUInt32(buf); err != nil {
		return fmt.Errorf("error while decoding time of sr: %w", err)
	} else {
		sr.NavigationTime = value
	}

	if value, err := common.ReadUInt32(buf); err != nil {
		return fmt.Errorf("error while decoding latitude of sr: %w", err)
	} else {
		sr.Latitude = value
	}

	if value, err := common.ReadUInt32(buf); err != nil {
		return fmt.Errorf("error while decoding longitude of sr: %w", err)
	} else {
		sr.Longitude = value
	}

	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while reading flags of sr: %w", err)
	}

	sr.AltitudeFieldExists = common.Bit(flags, 7)
	sr.LongitudeHemisphere = common.Bit(flags, 6)
	sr.LatitudeHemisphere = common.Bit(flags, 5)
	sr.Moving = common.Bit(flags, 4)
	sr.BlackBox = common.Bit(flags, 3)
	sr.CoordinateSystem = common.Bit(flags, 2)
	sr.Fixture = common.Bit(flags, 1)
	sr.Validness = common.Bit(flags, 0)

	if sr.Speed, err = common.ReadUInt16(buf); err != nil {
		return fmt.Errorf("error while decoding speed of sr: %w", err)
	}

	sr.AltitudeSign = sr.Speed&ALTSMask == ALTSMask
	sr.DirectionHighestBit = sr.Speed&DIRHMask == DIRHMask

	sr.Speed &= SpeedMask

	if sr.Direction, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding direction of sr: %w", err)
	}

	if sr.Odometer, err = common.ReadUInt24(buf); err != nil {
		return fmt.Errorf("error while decoding odometer of sr: %w", err)
	}

	if sr.DigitalInputs, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding digital inputs of sr: %w", err)
	}

	if sr.Source, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding source of sr: %w", err)
	}

	if sr.Altitude, err = common.ReadUInt24(buf); err != nil {
		return fmt.Errorf("error while decoding altitude of sr: %w", err)
	}

	if sr.SourceData, err = common.ReadUInt16(buf); err != nil {
		return fmt.Errorf("error while decoding source data of sr: %w", err)
	}

	return nil
}

func (sr *SrPosData) Encode() ([]byte, error) {
	var (
		flags        byte
		err          error
		uint32Buffer = make([]byte, 4)
	)

	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, sr.NavigationTime); err != nil {
		return nil, fmt.Errorf("error while encoding navigation time of sr: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, sr.Latitude); err != nil {
		return nil, fmt.Errorf("error while encoding latitude of sr: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, sr.Longitude); err != nil {
		return nil, fmt.Errorf("error while encoding longitude of sr: %w", err)
	}

	common.SetBit(&flags, 7, sr.AltitudeFieldExists)
	common.SetBit(&flags, 6, sr.LongitudeHemisphere)
	common.SetBit(&flags, 5, sr.LatitudeHemisphere)
	common.SetBit(&flags, 4, sr.Moving)
	common.SetBit(&flags, 3, sr.BlackBox)
	common.SetBit(&flags, 2, sr.CoordinateSystem)
	common.SetBit(&flags, 1, sr.Fixture)
	common.SetBit(&flags, 0, sr.Validness)

	if err = buf.WriteByte(flags); err != nil {
		return nil, fmt.Errorf("error while encoding flags of sr: %w", err)
	}

	packetSpeed := sr.Speed & SpeedMask

	if sr.AltitudeSign {
		packetSpeed |= ALTSMask
	}

	sr.DirectionHighestBit = (sr.Direction & DirMask) != 0

	if sr.Direction&DirMask != 0 {
		packetSpeed |= DIRHMask
	}

	if err = binary.Write(buf, binary.LittleEndian, packetSpeed); err != nil {
		return nil, fmt.Errorf("error while encoding speed of sr: %w", err)
	}

	if err = buf.WriteByte(sr.Direction); err != nil {
		return nil, fmt.Errorf("error while encoding direction of sr: %w", err)
	}

	binary.LittleEndian.PutUint32(uint32Buffer, sr.Odometer)

	if err = binary.Write(buf, binary.LittleEndian, uint32Buffer[:3]); err != nil {
		return nil, fmt.Errorf("error while encoding odometer of sr: %w", err)
	}

	if err = buf.WriteByte(sr.DigitalInputs); err != nil {
		return nil, fmt.Errorf("error while encoding digital inputs of sr: %w", err)
	}

	if err = buf.WriteByte(sr.Source); err != nil {
		return nil, fmt.Errorf("error while encoding source of sr: %w", err)
	}

	binary.LittleEndian.PutUint32(uint32Buffer, sr.Altitude)

	if err = binary.Write(buf, binary.LittleEndian, uint32Buffer[:3]); err != nil {
		return nil, fmt.Errorf("error while encoding altitude of sr: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, sr.SourceData); err != nil {
		return nil, fmt.Errorf("error while encoding source data of sr: %w", err)
	}

	return buf.Bytes(), nil
}

func (sr *SrPosData) Length() uint16 {
	if srBytes, err := sr.Encode(); err != nil {
		return 0
	} else {
		return uint16(len(srBytes))
	}
}
