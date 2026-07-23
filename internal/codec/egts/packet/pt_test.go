package packet

import (
	"reflect"
	"testing"
	"time"

	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/codec/egts/record"
	"github.com/agrotom/emulator/internal/codec/egts/teledata"
)

var timeOffset = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestPtAppData(t *testing.T) {
	pt := PtAppData{
		ServiceDataRecord{
			RecordNumber:             5,
			SourceServiceOnDevice:    true,
			RecipientServiceOnDevice: true,
			Group:                    false,
			RecordProcessingPriority: common.High,
			TimeFieldExists:          false,
			EventIDFieldExists:       false,
			ObjectIDFieldExists:      false,
			SourceServiceType:        0,
			RecipientServiceType:     0,
			RecordDataSet:            record.RecordDataSet{},
		},
	}

	_, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestFullPtAppData(t *testing.T) {
	pt := PtAppData{
		ServiceDataRecord{
			RecordNumber:             5,
			SourceServiceOnDevice:    true,
			RecipientServiceOnDevice: true,
			Group:                    false,
			RecordProcessingPriority: common.High,
			TimeFieldExists:          true,
			EventIDFieldExists:       true,
			ObjectIDFieldExists:      true,
			Time:                     uint32(time.Now().Unix()),
			EventIdentifier:          45,
			ObjectIdentifier:         31,
			SourceServiceType:        0,
			RecipientServiceType:     0,
			RecordDataSet:            record.RecordDataSet{},
		},
	}

	_, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestPtAppDataResponse(t *testing.T) {
	pt := PtResponse{
		ResponsePacketID: 5,
		ProcessingResult: 2,
		PtAppData: PtAppData{
			ServiceDataRecord{
				RecordNumber:             5,
				SourceServiceOnDevice:    true,
				RecipientServiceOnDevice: true,
				Group:                    false,
				RecordProcessingPriority: common.High,
				TimeFieldExists:          false,
				EventIDFieldExists:       false,
				ObjectIDFieldExists:      false,
				SourceServiceType:        0,
				RecipientServiceType:     0,
				RecordDataSet:            record.RecordDataSet{},
			},
		},
	}

	_, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestPtSignedAppData(t *testing.T) {
	pt := PtSignedAppData{
		SignatureLength: 4,
		SignatureData:   []byte{0x01, 0x02, 0x03, 0x04},
		PtAppData: PtAppData{
			ServiceDataRecord{
				RecordNumber:             5,
				SourceServiceOnDevice:    true,
				RecipientServiceOnDevice: true,
				Group:                    false,
				RecordProcessingPriority: common.High,
				TimeFieldExists:          false,
				EventIDFieldExists:       false,
				ObjectIDFieldExists:      false,
				SourceServiceType:        0,
				RecipientServiceType:     0,
				RecordDataSet:            record.RecordDataSet{},
			},
		},
	}

	_, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestPtAppDecoding(t *testing.T) {
	pt := PtAppData{
		ServiceDataRecord{
			RecordNumber:             5,
			SourceServiceOnDevice:    true,
			RecipientServiceOnDevice: true,
			Group:                    false,
			RecordProcessingPriority: common.High,
			TimeFieldExists:          false,
			EventIDFieldExists:       false,
			ObjectIDFieldExists:      false,
			SourceServiceType:        0,
			RecipientServiceType:     0,
			RecordDataSet:            record.RecordDataSet{},
		},
	}

	bytes, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	decodedPt := PtAppData{}
	err = decodedPt.Decode(bytes)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if !reflect.DeepEqual(pt, decodedPt) {
		t.Errorf("Encoded pt and decoded one have not been coincided")
	}
}

func TestFullPtAppDecoding(t *testing.T) {
	pt := PtAppData{
		ServiceDataRecord{
			RecordNumber:             5,
			SourceServiceOnDevice:    true,
			RecipientServiceOnDevice: true,
			Group:                    false,
			RecordProcessingPriority: common.High,
			TimeFieldExists:          true,
			EventIDFieldExists:       true,
			ObjectIDFieldExists:      true,
			Time:                     uint32(time.Now().Unix()),
			EventIdentifier:          45,
			ObjectIdentifier:         31,
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
						Moving:              false,
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
						Moving:              false,
						BlackBox:            true,
						Fixture:             false,
						CoordinateSystem:    false,
						Validness:           true,
						Speed:               332,
						AltitudeSign:        false,
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
	}

	bytes, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	decodedPt := PtAppData{}
	err = decodedPt.Decode(bytes)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	t.Logf("orig:    %+v", pt)
	t.Logf("decoded: %+v", decodedPt)

	if !reflect.DeepEqual(pt, decodedPt) {
		t.Errorf("Encoded pt and decoded one have not been coincided")
	}
}

func TestPtResponseDecoding(t *testing.T) {
	pt := PtResponse{
		ResponsePacketID: 5,
		ProcessingResult: 2,
		PtAppData: PtAppData{
			ServiceDataRecord{
				RecordNumber:             5,
				SourceServiceOnDevice:    true,
				RecipientServiceOnDevice: true,
				Group:                    true,
				RecordProcessingPriority: common.Low,
				TimeFieldExists:          false,
				EventIDFieldExists:       false,
				ObjectIDFieldExists:      false,
				SourceServiceType:        1,
				RecipientServiceType:     1,
				RecordDataSet:            record.RecordDataSet{},
			},
		},
	}

	bytes, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	decodedPt := PtResponse{}
	err = decodedPt.Decode(bytes)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	t.Logf("orig:    %+v", pt)
	t.Logf("decoded: %+v", decodedPt)

	if !reflect.DeepEqual(pt, decodedPt) {
		t.Errorf("Encoded pt response and decoded one have not been coincided")
	}
}

func TestPtSignedAppDataDecoding(t *testing.T) {
	pt := PtSignedAppData{
		SignatureLength: 4,
		SignatureData:   []byte{0x01, 0x02, 0x03, 0x04},
		PtAppData: PtAppData{
			ServiceDataRecord{
				RecordNumber:             5,
				SourceServiceOnDevice:    true,
				RecipientServiceOnDevice: true,
				Group:                    true,
				RecordProcessingPriority: common.Low,
				TimeFieldExists:          false,
				EventIDFieldExists:       false,
				ObjectIDFieldExists:      false,
				SourceServiceType:        1,
				RecipientServiceType:     1,
				RecordDataSet:            record.RecordDataSet{},
			},
		},
	}

	bytes, err := pt.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	decodedPt := PtSignedAppData{}
	err = decodedPt.Decode(bytes)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	t.Logf("orig:    %+v", pt)
	t.Logf("decoded: %+v", decodedPt)

	if !reflect.DeepEqual(pt, decodedPt) {
		t.Errorf("Encoded signed pt and decoded one have not been coincided")
	}
}
