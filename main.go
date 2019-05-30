package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/dakraid/skyrimSaveMaster/rgb"
	"github.com/dakraid/skyrimSaveMaster/tesvStruct"
	"github.com/dakraid/skyrimSaveMaster/typeReader"
)

const file = "DefaultSave.ess"

const (
	debug       = true
	printOffset = false
)

var saveGame tesvStruct.SaveFile

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readHeader(fileIn *os.File, startingOffset int64) int64 {
	var nextOffset = startingOffset

	saveGame.Header.Version, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.Header.SaveNumber, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.Header.PlayerName, nextOffset = typeReader.ReadWString(fileIn, nextOffset)
	saveGame.Header.PlayerLevel, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.Header.PlayerLocation, nextOffset = typeReader.ReadWString(fileIn, nextOffset)
	saveGame.Header.GameDate, nextOffset = typeReader.ReadWString(fileIn, nextOffset)
	saveGame.Header.PlayerRaceEditorId, nextOffset = typeReader.ReadWString(fileIn, nextOffset)
	saveGame.Header.PlayerSex, nextOffset = typeReader.ReadUInt16(fileIn, nextOffset)
	saveGame.Header.PlayerCurExp, nextOffset = typeReader.ReadFloat32(fileIn, nextOffset)
	saveGame.Header.PlayerLvlUpExp, nextOffset = typeReader.ReadFloat32(fileIn, nextOffset)
	saveGame.Header.Filetime, nextOffset = typeReader.ReadFiletime(fileIn, nextOffset)
	saveGame.Header.ShotWidth, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.Header.ShotHeight, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)

	return nextOffset
}

func readPlugins(fileIn *os.File, startingOffset int64) int64 {
	var nextOffset = startingOffset

	saveGame.PluginInfo.PluginCount, nextOffset = typeReader.ReadUInt8(fileIn, nextOffset)
	saveGame.PluginInfo.Plugins = make([]string, saveGame.PluginInfo.PluginCount)

	for i := 0; i < int(saveGame.PluginInfo.PluginCount); i++ {
		var temp string
		temp, nextOffset = typeReader.ReadWString(fileIn, nextOffset)
		saveGame.PluginInfo.Plugins[i] = temp
	}

	return nextOffset
}

func readFileLocationTable(fileIn *os.File, startingOffset int64) int64 {
	var nextOffset = startingOffset

	saveGame.FileLocationTable.FormIDArrayCountOffset, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.UnknownTable3Offset, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.GlobalDataTable1Offset, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.GlobalDataTable2Offset, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.ChangeFormsOffset, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.GlobalDataTable3Offset, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.GlobalDataTable1Count, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.GlobalDataTable2Count, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.GlobalDataTable3Count, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
	saveGame.FileLocationTable.ChangeFormCount, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)

	for i := 0; i < len(saveGame.FileLocationTable.Unused); i++ {
		var temp uint32
		temp, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)
		saveGame.FileLocationTable.Unused[i] = temp
	}

	return nextOffset
}

func readSkyrimLE(fileIn *os.File, startingOffset int64) {
	var nextOffset = startingOffset

	saveGame.HeaderSize, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)

	nextOffset = readHeader(fileIn, nextOffset)

	saveGame.Screenshot = rgb.NewImage(image.Rect(0, 0, int(saveGame.Header.ShotWidth), int(saveGame.Header.ShotHeight)))
	saveGame.Screenshot.Pix, nextOffset = typeReader.ReadScreenshot(fileIn, nextOffset, saveGame.Header.ShotWidth, saveGame.Header.ShotHeight)

	saveGame.FormVersion, nextOffset = typeReader.ReadUInt8(fileIn, nextOffset)
	saveGame.PluginInfoSize, nextOffset = typeReader.ReadUInt32(fileIn, nextOffset)

	nextOffset = readPlugins(fileIn, nextOffset)
	nextOffset = readFileLocationTable(fileIn, nextOffset)
}

func exportSave(filename string, source tesvStruct.SaveFile) {
	txtOut, err := os.Create(filename + ".txt")
	check(err)
	defer txtOut.Close()

	_, err = txtOut.WriteString(spew.Sdump(source))
	check(err)

	imgOut, err := os.Create(filename + ".png")
	check(err)
	defer imgOut.Close()

	err = png.Encode(imgOut, saveGame.Screenshot)
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

	magicCheck, err := typeReader.CheckMagic(f, saveGame.Magic)
	check(err)

	if magicCheck {
		readSkyrimLE(f, int64(13))
		exportSave("TES5_"+strings.TrimSuffix(filename, filepath.Ext(filename))+"_EX", saveGame)
	}
}
