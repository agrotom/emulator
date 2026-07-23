package record

import (
	"testing"
	"time"

	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/codec/egts/teledata"
)

var timeOffset = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestRecordEncoding(t *testing.T) {
	rd := RecordDataSet{
		RecordData{
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
	}

	_, err := rd.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestRecordDecoding(t *testing.T) {
	rd := RecordDataSet{
		RecordData{
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
	}

	rdBytes, err := rd.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	decodedRd := RecordDataSet{}
	err = decodedRd.Decode(rdBytes)

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestMultipleRecordDecoding(t *testing.T) {
	rd := RecordDataSet{
		RecordData{
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
		RecordData{
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
	}

	rdBytes, err := rd.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	decodedRd := RecordDataSet{}
	err = decodedRd.Decode(rdBytes)

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}
