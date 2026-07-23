package egts

import (
	"reflect"
	"testing"
	"time"

	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/codec/egts/packet"
	"github.com/agrotom/emulator/internal/codec/egts/record"
	"github.com/agrotom/emulator/internal/codec/egts/teledata"
	"github.com/google/go-cmp/cmp"
)

var timeOffset = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestEGTSFrameEncoding(t *testing.T) {
	frame := EGTSFrame{
		ProtocolVersion:   1,
		SecurityKeyID:     5,
		Prefix:            common.DefaultPrefix,
		Route:             false,
		Encryption:        common.NoEnc,
		Compressed:        false,
		Priority:          common.High,
		HeaderEncoding:    0,
		PacketIdentifier:  5,
		PacketType:        common.EGTSPTAppData,
		PeerAddress:       453,
		RecipientAddress:  554,
		TimeToLive:        3,
		ServicesFrameData: &packet.PtAppData{},
	}

	_, err := frame.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestEGTSFrameDecoding(t *testing.T) {
	frame := EGTSFrame{
		ProtocolVersion:   1,
		SecurityKeyID:     5,
		Prefix:            common.DefaultPrefix,
		Route:             true,
		Encryption:        common.NoEnc,
		Compressed:        false,
		Priority:          common.High,
		HeaderEncoding:    0,
		PacketIdentifier:  5,
		PacketType:        common.EGTSPTAppData,
		PeerAddress:       453,
		RecipientAddress:  554,
		TimeToLive:        3,
		ServicesFrameData: &packet.PtAppData{},
	}

	bytes, err := frame.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	decodedFrame := EGTSFrame{}
	err = decodedFrame.Decode(bytes)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if diff := cmp.Diff(frame, decodedFrame); diff != "" {
		t.Fatal(diff)
	}

	if !reflect.DeepEqual(frame, decodedFrame) {
		t.Errorf("Encoded frame and decoded one have not been coincided")
	}
}

func TestEGTSFrameFullDecoding(t *testing.T) {
	frame := EGTSFrame{
		ProtocolVersion:  1,
		SecurityKeyID:    5,
		Prefix:           common.DefaultPrefix,
		Route:            true,
		Encryption:       common.NoEnc,
		Compressed:       false,
		Priority:         common.High,
		HeaderEncoding:   0,
		PacketIdentifier: 5,
		PacketType:       common.EGTSPTAppData,
		PeerAddress:      453,
		RecipientAddress: 554,
		TimeToLive:       3,
		ServicesFrameData: &packet.PtAppData{
			packet.ServiceDataRecord{
				RecordNumber:             5,
				SourceServiceOnDevice:    true,
				RecipientServiceOnDevice: true,
				Group:                    false,
				RecordProcessingPriority: common.High,
				TimeFieldExists:          true,
				EventIDFieldExists:       true,
				ObjectIDFieldExists:      true,
				Time:                     uint32(time.Now().Unix()),
				EventIdentifier:          34,
				ObjectIdentifier:         45,
				SourceServiceType:        0,
				RecipientServiceType:     0,
				RecordDataSet: record.RecordDataSet{
					record.RecordData{
						SubRecordType: common.EgtsSrPosData,
						SubrecordData: &teledata.SrPosData{
							NavigationTime:      uint32(time.Since(timeOffset).Seconds()),
							Latitude:            434,
							Longitude:           353,
							AltitudeFieldExists: true,
							LongitudeHemisphere: true,
							LatitudeHemisphere:  true,
							Moving:              true,
							BlackBox:            true,
							Fixture:             false,
							CoordinateSystem:    false,
							Validness:           true,
							Speed:               345,
							AltitudeSign:        true,
							Direction:           242,
							Odometer:            342,
							DigitalInputs:       0,
							Source:              4,
							Altitude:            34,
							SourceData:          32,
						},
					},
					record.RecordData{
						SubRecordType: common.EgtsSrPosData,
						SubrecordData: &teledata.SrPosData{
							NavigationTime:      uint32(time.Since(timeOffset).Seconds()),
							Latitude:            222,
							Longitude:           132,
							AltitudeFieldExists: true,
							LongitudeHemisphere: true,
							LatitudeHemisphere:  true,
							Moving:              true,
							BlackBox:            true,
							Fixture:             false,
							CoordinateSystem:    false,
							Validness:           true,
							Speed:               332,
							AltitudeSign:        true,
							Direction:           255,
							Odometer:            3,
							DigitalInputs:       4,
							Source:              1,
							Altitude:            76,
							SourceData:          55,
						},
					},
				},
			},
		},
	}

	frameBytes, err := frame.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	decodedFrame := EGTSFrame{}
	err = decodedFrame.Decode(frameBytes)

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	if diff := cmp.Diff(frame, decodedFrame); diff != "" {
		t.Fatal(diff)
	}

	if !reflect.DeepEqual(frame, decodedFrame) {
		t.Errorf("Encoded frame and decoded one have not been coincided")
	}
}
