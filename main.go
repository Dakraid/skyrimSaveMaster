package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

const file = "TestSave.ess"
const debug = false

const (
	magicSize    = 13
	uint16Size   = 2
	uint32Size   = 4
	float32Size  = 4
	filetimeSize = 8
)

type saveHeader struct {
	headerSize         uint32
	version            uint32
	saveNumber         uint32
	playerName         string
	playerLevel        uint32
	playerLocation     string
	gameDate           string
	playerRaceEditorId string
	playerSex          uint16
	playerCurExp       float32
	playerLvlUpExp     float32
	filetime           time.Time
	shotWidth          uint32
	shotHeight         uint32
}

var saveGame saveHeader

var delta = time.Date(1970-369, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readWString(fileIn *os.File, offset int64) (string, int64) {
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

func readFiletime(fileIn *os.File, offset int64) (time.Time, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, filetimeSize)
	_, err = io.ReadAtLeast(fileIn, bytes, filetimeSize)
	check(err)

	t := binary.LittleEndian.Uint64(bytes)

	return time.Unix(0, int64(t)*100+delta), offset + filetimeSize
}

func readUInt16(fileIn *os.File, offset int64) (uint16, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, uint16Size)
	_, err = io.ReadAtLeast(fileIn, bytes, uint16Size)
	check(err)

	return binary.LittleEndian.Uint16(bytes), offset + uint16Size
}

func readUInt32(fileIn *os.File, offset int64) (uint32, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, uint32Size)
	_, err = io.ReadAtLeast(fileIn, bytes, uint32Size)
	check(err)

	return binary.LittleEndian.Uint32(bytes), offset + uint32Size
}

func readFloat32(fileIn *os.File, offset int64) (float32, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, float32Size)
	_, err = io.ReadAtLeast(fileIn, bytes, float32Size)
	check(err)

	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)

	return float, offset + float32Size
}

func checkMagic(fileIn *os.File) (bool, error) {
	magic := make([]byte, magicSize)
	_, err := fileIn.Read(magic)

	if err != nil {
		return false, err
	}

	if string(magic) == "TESV_SAVEGAME" {
		return true, nil
	} else {
		return false, errors.New("magic did not match expected value")
	}
}

func readHeader(fileIn *os.File) int64 {
	var nextOffset = int64(13)

	saveGame.headerSize, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.version, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.saveNumber, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.playerName, nextOffset = readWString(fileIn, nextOffset)
	saveGame.playerLevel, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.playerLocation, nextOffset = readWString(fileIn, nextOffset)
	saveGame.gameDate, nextOffset = readWString(fileIn, nextOffset)
	saveGame.playerRaceEditorId, nextOffset = readWString(fileIn, nextOffset)
	saveGame.playerSex, nextOffset = readUInt16(fileIn, nextOffset)
	saveGame.playerCurExp, nextOffset = readFloat32(fileIn, nextOffset)
	saveGame.playerLvlUpExp, nextOffset = readFloat32(fileIn, nextOffset)
	saveGame.filetime, nextOffset = readFiletime(fileIn, nextOffset)
	saveGame.shotWidth, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.shotHeight, nextOffset = readUInt32(fileIn, nextOffset)

	if debug {
		fmt.Printf("headerSize Offset: %d\n", nextOffset)
		fmt.Printf("version Offset: %d\n", nextOffset)
		fmt.Printf("saveNumber Offset: %d\n", nextOffset)
		fmt.Printf("playerName: %d\n", nextOffset)
		fmt.Printf("playerLevel Offset: %d\n", nextOffset)
		fmt.Printf("playerLocation Offset: %d\n", nextOffset)
		fmt.Printf("gameDate Offset: %d\n", nextOffset)
		fmt.Printf("playerRaceEditorId Offset: %d\n", nextOffset)
		fmt.Printf("playerSex Offset: %d\n", nextOffset)
		fmt.Printf("playerCurExp Offset: %d\n", nextOffset)
		fmt.Printf("playerLvlUpExp Offset: %d\n", nextOffset)
		fmt.Printf("filetime Offset: %d\n", nextOffset)
		fmt.Printf("shotWidth Offset: %d\n", nextOffset)
		fmt.Printf("shotHeight Offset: %d\n", nextOffset)
	}

	return int64(13 + saveGame.headerSize)
}

func main() {

	f, err := os.Open(file)
	check(err)
	defer f.Close()

	magicCheck, err := checkMagic(f)
	check(err)

	if magicCheck {
		readHeader(f)
	}

	fmt.Printf("%+v\n", saveGame)

	cet, _ := time.LoadLocation("CET")
	fmt.Println(saveGame.filetime.In(cet))
}
