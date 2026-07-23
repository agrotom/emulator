package wialonips

import (
	"fmt"
	"strings"
	"time"
)

type DataPacket struct {
	Date                          *time.Time
	Time                          *time.Time
	Latitude                      *float64
	Longitude                     *float64
	Speed                         *int
	Course                        *int
	Height                        *int
	Sats                          *int
	HorizontalDilutionOfPrecision *float64
	Inputs                        *int
	Outputs                       *int
	AnalogInputs                  *string
	IButton                       *string
	Options                       []Option
}

func (dp *DataPacket) Encode() ([]byte, error) {
	fstr := "#D#"

	WriteDate(&fstr, dp.Date)
	WriteTime(&fstr, dp.Time)

	if dp.Latitude != nil {
		lat, latDir := FormatLatCoord(*dp.Latitude)
		fstr += fmt.Sprintf("%s;%s;", lat, latDir)
	} else {
		fstr += "NA;NA;"
	}

	if dp.Longitude != nil {
		lon, lonDir := FormatLonCoord(*dp.Longitude)
		fstr += fmt.Sprintf("%s;%s;", lon, lonDir)
	} else {
		fstr += "NA;NA;"
	}

	WriteInt(&fstr, dp.Speed)
	WriteInt(&fstr, dp.Course)
	WriteInt(&fstr, dp.Height)
	WriteInt(&fstr, dp.Sats)
	WriteFloat64(&fstr, dp.HorizontalDilutionOfPrecision)
	WriteInt(&fstr, dp.Inputs)
	WriteInt(&fstr, dp.Outputs)
	WriteString(&fstr, dp.AnalogInputs)
	WriteStringNA(&fstr, dp.IButton)

	for _, opt := range dp.Options {
		fstr += opt.Format() + ";"
	}

	if len(dp.Options) > 0 {
		fstr = strings.TrimRight(fstr, ";")
	}

	return []byte(fstr + "\r\n"), nil
}

func (dp *DataPacket) ValidateResponse(resp []byte) error {
	if string(resp) != SuccessAD {
		return fmt.Errorf("%w: %s", ErrBadResponse, string(resp))
	}

	return nil
}
