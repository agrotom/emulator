package record

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/codec/egts/teledata"
)

type RecordData struct {
	SubRecordType   common.EgtsService
	SubRecordLength uint16
	SubrecordData   common.Binary
}

type RecordDataSet []RecordData

func (ds *RecordDataSet) Decode(data []byte) error {
	var (
		err               error
		subRecordTypeByte byte
	)

	buf := bytes.NewReader(data)

	for buf.Len() > 0 {
		rd := RecordData{}

		subRecordTypeByte, err = buf.ReadByte()

		if err != nil {
			return fmt.Errorf("error while decoding sub record type: %w", err)
		}

		rd.SubRecordType = common.EgtsService(subRecordTypeByte)

		if value, err := common.ReadUInt16(buf); err != nil {
			return fmt.Errorf("error while decoding sub record length: %w", err)
		} else {
			rd.SubRecordLength = value
		}

		if rd.SubrecordData, err = teledata.Factory(rd.SubRecordType); err != nil {
			return fmt.Errorf("error while decoding sub record data: %w", err)
		}

		recordBuffer := make([]byte, rd.SubRecordLength)
		if _, err = buf.Read(recordBuffer); err != nil {
			return fmt.Errorf("error while reading sub record data: %w", err)
		}

		rd.SubrecordData.Decode(recordBuffer)
		*ds = append(*ds, rd)
	}

	return nil
}

func (ds *RecordDataSet) Encode() ([]byte, error) {
	var (
		err error
	)

	buf := new(bytes.Buffer)

	for i := range *ds {
		rd := &(*ds)[i]

		if err = buf.WriteByte(byte(rd.SubRecordType)); err != nil {
			return nil, fmt.Errorf("error while encoding sub record type: %w", err)
		}

		if rd.SubRecordLength == 0 {
			rd.SubRecordLength = rd.SubrecordData.Length()
		}

		if err = binary.Write(buf, binary.LittleEndian, rd.SubRecordLength); err != nil {
			return nil, fmt.Errorf("error while encoding sub record length: %w", err)
		}

		srdBytes, err := rd.SubrecordData.Encode()

		if err != nil {
			return nil, fmt.Errorf("error while encoding sub record data: %w", err)
		}

		if _, err = buf.Write(srdBytes); err != nil {
			return nil, fmt.Errorf("error while encoding sub record data: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func (rd *RecordDataSet) Length() uint16 {
	if ptBytes, err := rd.Encode(); err != nil {
		return 0
	} else {
		return uint16(len(ptBytes))
	}
}
