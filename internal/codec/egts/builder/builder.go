package builder

import (
	"github.com/agrotom/emulator/internal/codec/egts"
	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/codec/egts/packet"
	"github.com/agrotom/emulator/internal/codec/egts/record"
)

type FrameBuilder struct {
	srds map[common.EgtsServiceType]*packet.ServiceDataRecord
}

func NewFrameBuilder() *FrameBuilder {
	builder := &FrameBuilder{}

	builder.srds = make(map[common.EgtsServiceType]*packet.ServiceDataRecord, 2)

	builder.srds[common.EgtsAuth] = &packet.ServiceDataRecord{
		SourceServiceOnDevice:    true,
		RecipientServiceOnDevice: false,
		SourceServiceType:        common.EgtsAuth,
		RecipientServiceType:     common.EgtsAuth,
		RecordDataSet:            record.RecordDataSet{},
	}

	builder.srds[common.EgtsTeledata] = &packet.ServiceDataRecord{
		SourceServiceOnDevice:    true,
		RecipientServiceOnDevice: false,
		SourceServiceType:        common.EgtsTeledata,
		RecipientServiceType:     common.EgtsTeledata,
		RecordDataSet:            record.RecordDataSet{},
	}

	return builder
}

func (b *FrameBuilder) AddAuthRecord(rtype common.EgtsService, data common.Binary) *FrameBuilder {
	b.srds[common.EgtsAuth].RecordDataSet = append(b.srds[common.EgtsAuth].RecordDataSet, record.RecordData{
		SubRecordType: rtype,
		SubrecordData: data,
	})

	return b
}

func (b *FrameBuilder) AddTeledataRecord(rtype common.EgtsService, data common.Binary) *FrameBuilder {
	b.srds[common.EgtsTeledata].RecordDataSet = append(b.srds[common.EgtsTeledata].RecordDataSet, record.RecordData{
		SubRecordType: rtype,
		SubrecordData: data,
	})

	return b
}

func (b *FrameBuilder) Build() *egts.EGTSFrame {
	frame := &egts.EGTSFrame{}
	payload := &packet.PtAppData{}

	frame.ProtocolVersion = 0x01
	frame.PacketType = common.EGTSPTAppData
	frame.ServicesFrameData = payload

	var recordNumber uint16 = 0

	if len(b.srds[common.EgtsAuth].RecordDataSet) > 0 {
		b.srds[common.EgtsAuth].RecordNumber = recordNumber
		recordNumber++
		*payload = append(*payload, *b.srds[common.EgtsAuth])
	}

	if len(b.srds[common.EgtsTeledata].RecordDataSet) > 0 {
		b.srds[common.EgtsTeledata].RecordNumber = recordNumber
		recordNumber++
		*payload = append(*payload, *b.srds[common.EgtsTeledata])
	}

	return frame
}
