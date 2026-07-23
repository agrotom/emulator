package common

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

func ReadUInt16(r *bytes.Reader) (uint16, error) {
	uint16Buffer := make([]byte, 2)

	if _, err := r.Read(uint16Buffer); err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(uint16Buffer), nil
}

func ReadUInt24(r *bytes.Reader) (uint32, error) {
	uint24Buffer := make([]byte, 3)

	if _, err := r.Read(uint24Buffer); err != nil {
		return 0, err
	}

	uint24Buffer = append(uint24Buffer, 0x00)

	return binary.LittleEndian.Uint32(uint24Buffer), nil
}

func WriteUInt24(w *bytes.Buffer, data uint32) error {
	uint32Buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(uint32Buffer, data)

	if err := binary.Write(w, binary.LittleEndian, uint32Buffer[:3]); err != nil {
		return fmt.Errorf("error while encoding altitude of sr: %w", err)
	}

	return nil
}

func ReadUInt32(r *bytes.Reader) (uint32, error) {
	uint32Buffer := make([]byte, 4)

	if _, err := r.Read(uint32Buffer); err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(uint32Buffer), nil
}

func ReadStringBinary(r *bytes.Reader, size uint16) (string, error) {
	var err error
	buffer := make([]byte, size)

	if _, err = r.Read(buffer); err != nil {
		return "", err
	}

	decoded, err := charmap.Windows1251.NewDecoder().Bytes(buffer)

	if err != nil {
		return "", err
	}

	decoded = bytes.TrimRight(decoded, "\x00")

	return string(decoded), nil
}

func WriteStringBinary(w *bytes.Buffer, size uint16, data string) error {
	encoded, err := charmap.Windows1251.NewEncoder().Bytes([]byte(data))

	if err != nil {
		return err
	}

	buffer := make([]byte, size)
	copy(buffer, encoded)

	binary.Write(w, binary.LittleEndian, buffer)

	return nil
}
