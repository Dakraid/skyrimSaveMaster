package typeReader

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"
	"time"
)

const (
	magicSize    = 13
	uint8Size    = 1
	uint16Size   = 2
	uint32Size   = 4
	float32Size  = 4
	filetimeSize = 8
)

var delta = time.Date(1970-369, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckMagic(fileIn *os.File, target string) (bool, error) {
	magic := make([]byte, magicSize)
	_, err := fileIn.Read(magic)

	if err != nil {
		return false, err
	}

	target = string(magic)
	if target == "TESV_SAVEGAME" {
		return true, nil
	} else {
		return false, errors.New("magic did not match expected value")
	}
}

func ReadWString(fileIn *os.File, offset int64) (string, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	sizeBytes := make([]byte, uint16Size)
	_, err = io.ReadAtLeast(fileIn, sizeBytes, uint16Size)
	check(err)

	wstringSize := binary.LittleEndian.Uint16(sizeBytes)

	_, err = fileIn.Seek(offset+2, 0)
	check(err)

	bytes := make([]byte, wstringSize)
	_, err = io.ReadAtLeast(fileIn, bytes, int(wstringSize))
	check(err)

	return string(bytes), offset + 2 + int64(wstringSize)
}

func ReadFiletime(fileIn *os.File, offset int64) (time.Time, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, filetimeSize)
	_, err = io.ReadAtLeast(fileIn, bytes, filetimeSize)
	check(err)

	t := binary.LittleEndian.Uint64(bytes)

	return time.Unix(0, int64(t)*100+delta), offset + filetimeSize
}

func ReadUInt8(fileIn *os.File, offset int64) (uint8, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, uint8Size)
	_, err = io.ReadAtLeast(fileIn, bytes, uint8Size)
	check(err)

	return bytes[0], offset + uint8Size
}

func ReadUInt16(fileIn *os.File, offset int64) (uint16, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, uint16Size)
	_, err = io.ReadAtLeast(fileIn, bytes, uint16Size)
	check(err)

	return binary.LittleEndian.Uint16(bytes), offset + uint16Size
}

func ReadUInt32(fileIn *os.File, offset int64) (uint32, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, uint32Size)
	_, err = io.ReadAtLeast(fileIn, bytes, uint32Size)
	check(err)

	return binary.LittleEndian.Uint32(bytes), offset + uint32Size
}

func ReadFloat32(fileIn *os.File, offset int64) (float32, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, float32Size)
	_, err = io.ReadAtLeast(fileIn, bytes, float32Size)
	check(err)

	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)

	return float, offset + float32Size
}

func ReadScreenshot(fileIn *os.File, startingOffset int64, shotWidth, shotHeight uint32) ([]uint8, int64) {
	var nextOffset = startingOffset

	arraySize := 3 * shotWidth * shotHeight

	_, err := fileIn.Seek(nextOffset, 0)
	check(err)

	screenshotData := make([]uint8, arraySize)
	_, err = io.ReadAtLeast(fileIn, screenshotData, int(arraySize))
	check(err)

	return screenshotData, nextOffset + int64(arraySize)
}