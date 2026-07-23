package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/codec/egts/record"
)

type ServiceDataRecord struct {
	RecordLength uint16
	RecordNumber uint16

	SourceServiceOnDevice    bool
	RecipientServiceOnDevice bool
	Group                    bool
	RecordProcessingPriority common.Priority
	TimeFieldExists          bool
	EventIDFieldExists       bool
	ObjectIDFieldExists      bool

	ObjectIdentifier uint32
	EventIdentifier  uint32
	Time             uint32

	SourceServiceType    common.EgtsServiceType
	RecipientServiceType common.EgtsServiceType
	record.RecordDataSet
}

type PtAppData []ServiceDataRecord

func (pt *PtAppData) Decode(data []byte) error {
	var (
		flags        byte
		err          error
		uint16Buffer = make([]byte, 2)
		uintBuffer   = make([]byte, 4)
	)

	buf := bytes.NewReader(data)

	for buf.Len() > 0 {
		sdr := ServiceDataRecord{}

		if _, err = buf.Read(uint16Buffer); err != nil {
			return fmt.Errorf("error while decoding record length of sdr: %w", err)
		}

		sdr.RecordLength = binary.LittleEndian.Uint16(uint16Buffer)

		if _, err = buf.Read(uint16Buffer); err != nil {
			return fmt.Errorf("error while decoding record number of sdr: %w", err)
		}

		sdr.RecordNumber = binary.LittleEndian.Uint16(uint16Buffer)

		if flags, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("error while decoding flags of sdr: %w", err)
		}

		sdr.SourceServiceOnDevice = common.Bit(flags, 7)
		sdr.RecipientServiceOnDevice = common.Bit(flags, 6)
		sdr.Group = common.Bit(flags, 5)
		sdr.RecordProcessingPriority = common.Priority(common.Bits(flags, 3, 2))
		sdr.TimeFieldExists = common.Bit(flags, 2)
		sdr.EventIDFieldExists = common.Bit(flags, 1)
		sdr.ObjectIDFieldExists = common.Bit(flags, 0)

		if sdr.ObjectIDFieldExists {
			if _, err = buf.Read(uintBuffer); err != nil {
				return fmt.Errorf("error while decoding object id of sdr: %w", err)
			}

			sdr.ObjectIdentifier = binary.LittleEndian.Uint32(uintBuffer)
		}

		if sdr.EventIDFieldExists {
			if _, err = buf.Read(uintBuffer); err != nil {
				return fmt.Errorf("error while decoding event id of sdr: %w", err)
			}

			sdr.EventIdentifier = binary.LittleEndian.Uint32(uintBuffer)
		}

		if sdr.TimeFieldExists {
			if _, err = buf.Read(uintBuffer); err != nil {
				return fmt.Errorf("error while decoding time of sdr: %w", err)
			}

			sdr.Time = binary.LittleEndian.Uint32(uintBuffer)
		}

		if sstype, err := buf.ReadByte(); err != nil {
			return fmt.Errorf("error while decoding source service type of sdr: %w", err)
		} else {
			sdr.SourceServiceType = common.EgtsServiceType(sstype)
		}

		if rctype, err := buf.ReadByte(); err != nil {
			return fmt.Errorf("error while decoding source service type of sdr: %w", err)
		} else {
			sdr.RecipientServiceType = common.EgtsServiceType(rctype)
		}

		sdr.RecordDataSet = record.RecordDataSet{}

		if buf.Len() != 0 {
			recordData := make([]byte, sdr.RecordLength)
			if _, err := buf.Read(recordData); err != nil {
				return fmt.Errorf("error while reading sdr record: %w", err)
			}

			if err = sdr.RecordDataSet.Decode(recordData); err != nil {
				return fmt.Errorf("error while decoding sdr record: %w", err)
			}
		}

		*pt = append(*pt, sdr)
	}

	return nil
}

func (pt *PtAppData) Encode() ([]byte, error) {
	var (
		flags byte
	)

	buf := new(bytes.Buffer)

	for i := range *pt {
		sdr := &(*pt)[i]
		res, err := sdr.RecordDataSet.Encode()

		if err != nil {
			return nil, fmt.Errorf("error while encoding record of sdr: %w", err)
		}

		if sdr.RecordLength == 0 {
			sdr.RecordLength = uint16(len(res))
		}

		if err = binary.Write(buf, binary.LittleEndian, sdr.RecordLength); err != nil {
			return nil, fmt.Errorf("error while encoding record length of sdr: %w", err)
		}

		if err = binary.Write(buf, binary.LittleEndian, sdr.RecordNumber); err != nil {
			return nil, fmt.Errorf("error while encoding record length of sdr: %w", err)
		}

		common.SetBit(&flags, 7, sdr.SourceServiceOnDevice)
		common.SetBit(&flags, 6, sdr.RecipientServiceOnDevice)
		common.SetBit(&flags, 5, sdr.Group)
		common.SetBits(&flags, 3, 2, byte(sdr.RecordProcessingPriority))
		common.SetBit(&flags, 2, sdr.TimeFieldExists)
		common.SetBit(&flags, 1, sdr.EventIDFieldExists)
		common.SetBit(&flags, 0, sdr.ObjectIDFieldExists)

		if err = buf.WriteByte(byte(flags)); err != nil {
			return nil, fmt.Errorf("error while encoding flags: %w", err)
		}

		if sdr.ObjectIDFieldExists {
			if err = binary.Write(buf, binary.LittleEndian, sdr.ObjectIdentifier); err != nil {
				return nil, fmt.Errorf("error while encoding object identifier of sdr: %w", err)
			}
		}

		if sdr.EventIDFieldExists {
			if err = binary.Write(buf, binary.LittleEndian, sdr.EventIdentifier); err != nil {
				return nil, fmt.Errorf("error while encoding event identifier of sdr: %w", err)
			}
		}

		if sdr.TimeFieldExists {
			if err = binary.Write(buf, binary.LittleEndian, sdr.Time); err != nil {
				return nil, fmt.Errorf("error while encoding time of sdr: %w", err)
			}
		}

		if err = buf.WriteByte(byte(sdr.SourceServiceType)); err != nil {
			return nil, fmt.Errorf("error while encoding source service type of sdr: %w", err)
		}

		if err = buf.WriteByte(byte(sdr.RecipientServiceType)); err != nil {
			return nil, fmt.Errorf("error while encoding recipient service type of sdr: %w", err)
		}

		buf.Write(res)
	}

	return buf.Bytes(), nil
}

func (pt *PtAppData) Length() uint16 {
	if ptBytes, err := pt.Encode(); err != nil {
		return 0
	} else {
		return uint16(len(ptBytes))
	}
}
