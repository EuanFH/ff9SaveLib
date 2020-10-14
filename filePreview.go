package FF9Save

import (
	"bytes"
	"encoding/binary"
)
const FilePreviewSize = 965
const FilePreviewReservedSize = 1024

type FilePreview struct {
	IsPreviewCorrupted bool
	HasData bool
	Gil int64
	PlayDuration uint64
	WinType uint64 `json:"win_type"`
	Location String128
	CharacterInfoList [4]CharacterInfo
	Timestamp float64
	ReservedData [64]int32
}

type CharacterInfo struct {
	SerialID int32
	Level int32
	Name String128
}

func (fp *FilePreview) BinaryMarshaler() ([]byte, error) {
	buf := new(bytes.Buffer)
	filePreview := *fp
	if *fp == (FilePreview{}) {
		if err := binary.Write(buf, binary.LittleEndian, []byte{'N','O','N','E'}); err != nil{
			return nil, err
		}
	} else {
		if err := binary.Write(buf, binary.LittleEndian, []byte{'P','R','E','V'}); err != nil{
			return nil, err
		}
		filePreview.FixCharacterInfoForBinary()
	}
	filePreviewBuf := new(bytes.Buffer)
	if err := binary.Write(filePreviewBuf, binary.LittleEndian, filePreview); err != nil{
		return nil, err
	}
	//remove IsPreviewCorrupted its no included in binary
	if err := binary.Write(buf, binary.LittleEndian, filePreviewBuf.Bytes()[1:]); err != nil{
		return nil, err
	}
	return buf.Bytes(), nil
}
//in json the values are 0 for empty characterInfo but in binary they are -1
func (fp *FilePreview) FixCharacterInfoForBinary() {
	for i, _ := range fp.CharacterInfoList {
		if fp.CharacterInfoList[i].SerialID == 0 {
			fp.CharacterInfoList[i].SerialID = -1
		}
		if fp.CharacterInfoList[i].Level == 0 {
			fp.CharacterInfoList[i].Level = -1
		}
	}
}

func (fp *FilePreview) FixCharacterInfoFromBinary() {
	for i, _ := range fp.CharacterInfoList {
		if fp.CharacterInfoList[i].SerialID == -1 {
			fp.CharacterInfoList[i].SerialID = 0
		}
		if fp.CharacterInfoList[i].Level == -1 {
			fp.CharacterInfoList[i].Level = 0
		}
	}
}

func (fp *FilePreview)UnBinaryMarshaler(data []byte) error{
	//remove header set is corrupted to false
	data = append([]byte{0x00}, data[4:]...)
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, fp); err != nil{
		return err
	}
	fp.FixCharacterInfoFromBinary()
	return nil
}

