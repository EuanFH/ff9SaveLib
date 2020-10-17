package ff9Save

import (
	"bytes"
	"encoding/binary"
)

const MetaDataSize = 288 //this is including header
const MetaDataReservedSize = 320

type MetaData struct {
	SaveVersion float32 //1f
	DataSize int32 //they always plus 4 onto this to account for the NONE or SAVE Header rembere to remove 4 bytes from my file size count
	FileInfo FileInfo
	IsGameFinishFlag int32
	SelectedLanguage int32 //PrefsLanguage //will have to figure out how to use the enum convertion on binary unmarshal
	IsAutoLogin int8
	SystemAchievementStatuses byte //this is an array but don't think it matters its a size of 1
	ScreenRotation byte
	ReservedBuffer [249]byte
}

func NewMetaData() MetaData{
	return MetaData{
		SaveVersion: 1.0,
		DataSize: FileSize - 4,
		IsGameFinishFlag: 0, //zero is not finished one is finished
		IsAutoLogin: 0x00 , //temp will need to read this from the real game file. no idea what are valid values
		SystemAchievementStatuses: 0x00,
		ScreenRotation: 0x03, //assuming this is landscape for this now. will need to read a normal save file //seems to be screen rotation is zero on pc
	}
}

func (md *MetaData) BinaryMarshaler() ([]byte, error){
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, []byte{'S','A','V','E'}); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian,  md); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}