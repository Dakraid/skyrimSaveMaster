package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/dakraid/skyrimSaveMaster/rgb"
)

const file = "DefaultSave.ess"

const (
	debug       = true
	printOffset = false
)

const (
	magicSize    = 13
	uint8Size    = 1
	uint16Size   = 2
	uint32Size   = 4
	float32Size  = 4
	filetimeSize = 8
)

type saveFile struct {
	magic             string
	headerSize        uint32
	header            saveHeader
	screenshot        *rgb.Image
	formVersion       uint8
	pluginInfoSize    uint32
	pluginInfo        savePlugins
	fileLocationTable saveFileLocationTable
}

type saveHeader struct {
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

type savePlugins struct {
	pluginCount uint8
	plugins     []string
}

type saveFileLocationTable struct {
	formIDArrayCountOffset uint32
	changeFormsOffset      uint32
	globalDataTable1Offset uint32
	globalDataTable2Offset uint32
	globalDataTable3Offset uint32
	unknownTable3Offset    uint32
	changeFormCount        uint32
	globalDataTable1Count  uint32
	globalDataTable2Count  uint32
	globalDataTable3Count  uint32
	unused                 [15]uint32
}

var saveGame saveFile

var delta = time.Date(1970-369, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkMagic(fileIn *os.File) (bool, error) {
	magic := make([]byte, magicSize)
	_, err := fileIn.Read(magic)

	if err != nil {
		return false, err
	}

	saveGame.magic = string(magic)
	if string(magic) == "TESV_SAVEGAME" {
		return true, nil
	} else {
		return false, errors.New("magic did not match expected value")
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

func readUInt8(fileIn *os.File, offset int64) (uint8, int64) {
	_, err := fileIn.Seek(offset, 0)
	check(err)

	bytes := make([]byte, uint8Size)
	_, err = io.ReadAtLeast(fileIn, bytes, uint8Size)
	check(err)

	return bytes[0], offset + uint8Size
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

func readScreenshot(fileIn *os.File, startingOffset int64) int64 {
	var nextOffset = startingOffset

	arraySize := 3 * saveGame.header.shotWidth * saveGame.header.shotHeight

	_, err := fileIn.Seek(nextOffset, 0)
	check(err)

	screenshotData := make([]uint8, arraySize)
	_, err = io.ReadAtLeast(fileIn, screenshotData, int(arraySize))
	check(err)

	saveGame.screenshot = rgb.NewImage(image.Rect(0, 0, int(saveGame.header.shotWidth), int(saveGame.header.shotHeight)))
	saveGame.screenshot.Pix = screenshotData

	return nextOffset + int64(arraySize)
}

func readHeader(fileIn *os.File, startingOffset int64) int64 {
	var nextOffset = startingOffset

	saveGame.header.version, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.header.saveNumber, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.header.playerName, nextOffset = readWString(fileIn, nextOffset)
	saveGame.header.playerLevel, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.header.playerLocation, nextOffset = readWString(fileIn, nextOffset)
	saveGame.header.gameDate, nextOffset = readWString(fileIn, nextOffset)
	saveGame.header.playerRaceEditorId, nextOffset = readWString(fileIn, nextOffset)
	saveGame.header.playerSex, nextOffset = readUInt16(fileIn, nextOffset)
	saveGame.header.playerCurExp, nextOffset = readFloat32(fileIn, nextOffset)
	saveGame.header.playerLvlUpExp, nextOffset = readFloat32(fileIn, nextOffset)
	saveGame.header.filetime, nextOffset = readFiletime(fileIn, nextOffset)
	saveGame.header.shotWidth, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.header.shotHeight, nextOffset = readUInt32(fileIn, nextOffset)

	if debug && printOffset {
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

	return nextOffset
}

func readPlugins(fileIn *os.File, startingOffset int64) int64 {
	var nextOffset = startingOffset

	saveGame.pluginInfo.pluginCount, nextOffset = readUInt8(fileIn, nextOffset)
	saveGame.pluginInfo.plugins = make([]string, saveGame.pluginInfo.pluginCount)

	for i := 0; i < int(saveGame.pluginInfo.pluginCount); i++ {
		var temp string
		temp, nextOffset = readWString(fileIn, nextOffset)
		saveGame.pluginInfo.plugins[i] = temp
	}

	return nextOffset
}

func readFileLocationTable(fileIn *os.File, startingOffset int64) int64 {
	var nextOffset = startingOffset

	saveGame.fileLocationTable.formIDArrayCountOffset, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.unknownTable3Offset, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.globalDataTable1Offset, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.globalDataTable2Offset, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.changeFormsOffset, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.globalDataTable3Offset, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.globalDataTable1Count, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.globalDataTable2Count, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.globalDataTable3Count, nextOffset = readUInt32(fileIn, nextOffset)
	saveGame.fileLocationTable.changeFormCount, nextOffset = readUInt32(fileIn, nextOffset)

	for i := 0; i < len(saveGame.fileLocationTable.unused); i++ {
		var temp uint32
		temp, nextOffset = readUInt32(fileIn, nextOffset)
		saveGame.fileLocationTable.unused[i] = temp
	}

	return nextOffset
}

func readSkyrimLE(fileIn *os.File, startingOffset int64) {
	var nextOffset = startingOffset

	saveGame.headerSize, nextOffset = readUInt32(fileIn, nextOffset)

	nextOffset = readHeader(fileIn, nextOffset)
	nextOffset = readScreenshot(fileIn, nextOffset)

	saveGame.formVersion, nextOffset = readUInt8(fileIn, nextOffset)
	saveGame.pluginInfoSize, nextOffset = readUInt32(fileIn, nextOffset)

	nextOffset = readPlugins(fileIn, nextOffset)
	nextOffset = readFileLocationTable(fileIn, nextOffset)
}

func exportSave(filename string, source saveFile) {
	txtOut, err := os.Create(filename + ".txt")
	check(err)
	defer txtOut.Close()

	_, err = txtOut.WriteString(spew.Sdump(source))
	check(err)

	imgOut, err := os.Create(filename + ".png")
	check(err)
	defer imgOut.Close()

	err = png.Encode(imgOut, saveGame.screenshot)
	check(err)

	/*
		imgRes := resize.Resize(0, 720, saveGame.screenshot, resize.MitchellNetravali)

		imgOutB, err := os.Create(filename + "_big.png")
		check(err)
		defer imgOutB.Close()

		err = png.Encode(imgOutB, imgRes)
		check(err)
	*/
}

func main() {
	var filename = file

	cliArgs := os.Args[1:]

	if len(cliArgs) > 0 && cliArgs[0] != "" {
		filename = cliArgs[0]
	}

	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	magicCheck, err := checkMagic(f)
	check(err)

	if magicCheck {
		readSkyrimLE(f, int64(13))
		exportSave("TES5_"+strings.TrimSuffix(filename, filepath.Ext(filename))+"_EX", saveGame)
	}
}
