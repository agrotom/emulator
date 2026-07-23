package wialonips

import (
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	UnitID   string
	Host     string
	Port     string
	Password string
)

func TestMain(m *testing.M) {
	UnitID = os.Getenv("EMULATOR_TEST_UNIT_ID")
	Host = os.Getenv("EMULATOR_TEST_HOST")
	Port = os.Getenv("EMULATOR_TEST_PORT")
	Password = os.Getenv("EMULATOR_TEST_PASSWORD")

	os.Exit(m.Run())
}

func TestLoginPacket(t *testing.T) {
	packet := &LoginPacket{
		IMEI:     UnitID,
		Password: Port,
	}

	expectedResult := []byte("#L#" + UnitID + ";" + Port + "\r\n")

	encodedPacket, _ := packet.Encode()
	if !reflect.DeepEqual(encodedPacket, expectedResult) {
		t.Errorf("packet converted to ascii string and expected ascii string are not coincided")
	}
}

func TestDataPacket(t *testing.T) {
	nowDate, _ := time.Parse("020106", "170726")
	nowTime, _ := time.Parse("150405", "012900")

	lat := 51.536434
	lon := 46.023555
	speed := 5
	course := 7
	height := 1
	sats := 8
	hdop := 0.654
	inputs := 42
	outputs := 3

	packet := &DataPacket{
		Date:                          &nowDate,
		Time:                          &nowTime,
		Latitude:                      &lat,
		Longitude:                     &lon,
		Speed:                         &speed,
		Course:                        &course,
		Height:                        &height,
		Sats:                          &sats,
		HorizontalDilutionOfPrecision: &hdop,
		Inputs:                        &inputs,
		Outputs:                       &outputs,
		AnalogInputs:                  nil,
		IButton:                       nil,
	}

	expectedResult := []byte("#D#170726;012900;5132.1860;N;4601.4133;E;5;7;1;8;0.654;42;3;;NA;\r\n")

	encodedPacket, _ := packet.Encode()
	if !reflect.DeepEqual(encodedPacket, expectedResult) {
		t.Log(string(encodedPacket))
		t.Errorf("packet converted to ascii string and expected ascii string are not coincided")
	}
}
