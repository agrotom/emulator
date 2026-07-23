package wialon

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/agrotom/emulator/internal/codec/wialonips"
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

func createWialonAndDial() (*TCPClient, error) {
	var err error
	client, err := CreateWialonClient(&wialonips.WialonIPSReader{}, Host, Port)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err = client.Dial(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func TestWialonDial(t *testing.T) {
	if _, err := createWialonAndDial(); err != nil {
		t.Errorf("error while dialing: %s", err.Error())
		return
	}
}

func WialonLogin(t *testing.T) {
	var (
		err    error
		client *TCPClient
	)

	if client, err = createWialonAndDial(); err != nil {
		t.Errorf("error while dialing: %s", err.Error())
		return
	}

	expectedData := []byte("#AL#1\r\n")
	data, err := client.SendPacket(context.Background(), &wialonips.LoginPacket{
		IMEI:     UnitID,
		Password: Password,
	})

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("%s (%x) doesn't equal %s (%x)", string(data), data, string(expectedData), expectedData)
	}
}

func WialonSendPacket(t *testing.T) {
	var (
		err    error
		client *TCPClient
	)

	if client, err = createWialonAndDial(); err != nil {
		t.Errorf("error while dialing: %s", err.Error())
		return
	}

	expectedData := []byte("#AD#1\r\n")

	dateTime := time.Now()
	lat := 51.536434
	lon := 46.023555
	speed := 120
	course := 30
	sats := 1
	inputs := 0

	packet := &wialonips.DataPacket{
		Date:                          &dateTime,
		Time:                          &dateTime,
		Latitude:                      &lat,
		Longitude:                     &lon,
		Speed:                         &speed,
		Course:                        &course,
		Height:                        nil,
		Sats:                          &sats,
		HorizontalDilutionOfPrecision: nil,
		Inputs:                        &inputs,
	}

	encodedPacket, _ := packet.Encode()
	t.Logf("%s", string(encodedPacket))

	data, err := client.SendPacket(context.Background(), packet)

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("%s (%x) doesn't equal %s (%x)", string(data), data, string(expectedData), expectedData)
	}
}
