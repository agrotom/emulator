package teledata

import (
	"reflect"
	"testing"
	"time"
)

var timeOffset = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestSrEncoding(t *testing.T) {
	sr := SrPosData{
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
	}

	_, err := sr.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestSrDecoding(t *testing.T) {

	sr := SrPosData{
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
	}

	srBytes, err := sr.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	decodedSr := SrPosData{}
	err = decodedSr.Decode(srBytes)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	t.Logf("orig:    %+v", sr)
	t.Logf("decoded: %+v", decodedSr)

	if !reflect.DeepEqual(sr, decodedSr) {
		t.Errorf("Encoded sr and decoded one have not been coincided")
	}
}
