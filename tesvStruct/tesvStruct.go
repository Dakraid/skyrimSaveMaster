package tesvStruct

import (
	"time"

	"github.com/dakraid/skyrimSaveMaster/rgb"
)

type SaveFile struct {
	Magic             string
	HeaderSize        uint32
	Header            saveHeader
	Screenshot        *rgb.Image
	FormVersion       uint8
	PluginInfoSize    uint32
	PluginInfo        savePlugins
	FileLocationTable saveFileLocationTable
}

type saveHeader struct {
	Version            uint32
	SaveNumber         uint32
	PlayerName         string
	PlayerLevel        uint32
	PlayerLocation     string
	GameDate           string
	PlayerRaceEditorId string
	PlayerSex          uint16
	PlayerCurExp       float32
	PlayerLvlUpExp     float32
	Filetime           time.Time
	ShotWidth          uint32
	ShotHeight         uint32
}

type savePlugins struct {
	PluginCount uint8
	Plugins     []string
}

type saveFileLocationTable struct {
	FormIDArrayCountOffset uint32
	ChangeFormsOffset      uint32
	GlobalDataTable1Offset uint32
	GlobalDataTable2Offset uint32
	GlobalDataTable3Offset uint32
	UnknownTable3Offset    uint32
	ChangeFormCount        uint32
	GlobalDataTable1Count  uint32
	GlobalDataTable2Count  uint32
	GlobalDataTable3Count  uint32
	Unused                 [15]uint32
}
