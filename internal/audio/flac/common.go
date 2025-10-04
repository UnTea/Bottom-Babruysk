package flac

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrNegativeLength = errors.New("negative length")
)

// readBytes считывает определённое количество байт.
func readBytes(r io.Reader, n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrNegativeLength
	}

	buffer := make([]byte, n)
	_, err := io.ReadFull(r, buffer)

	return buffer, err
}

// readU16BE считывает uint16 в формате big-endian.
func readU16BE(r io.Reader) (uint16, error) {
	var buffer [2]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(buffer[:]), nil
}

// readU24BE считывает uint24 в формате big-endian.
func readU24BE(r io.Reader) (uint32, error) {
	var buffer [3]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, err
	}

	data0 := uint32(buffer[0]) << 16
	data1 := uint32(buffer[1]) << 8
	data2 := uint32(buffer[2])

	return data0 | data1 | data2, nil
}

// readU32BE считывает uint32 в формате big-endian.
func readU32BE(r io.Reader) (uint32, error) {
	var buffer [4]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(buffer[:]), nil
}

// readU32LE считывает uint32 в формате little-endian.
func readU32LE(r io.Reader) (uint32, error) {
	var buffer [4]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buffer[:]), nil
}

// readU64BE считывает uint64 в формате big-endian.
func readU64BE(r io.Reader) (uint64, error) {
	var buffer [8]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(buffer[:]), nil
}

// readLPStringLE считывает length-prefixed строку: сначала uint32 LE длины, затем байты строки.
func readLPStringLE(r io.Reader) (string, error) {
	n, err := readU32LE(r)
	if err != nil {
		return "", err
	}

	bytes, err := readBytes(r, int(n))
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
