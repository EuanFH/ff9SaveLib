package FF9Save

import (
	"bytes"
	"encoding/binary"
)
//size=961
//savecount=15

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

func (fp *FilePreview) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if *fp == (FilePreview{}) {
		if err := binary.Write(buf, binary.LittleEndian, []byte{'N','O','N','E'}); err != nil{
			return nil, err
		}
		//write zeros for len
		return buf.Bytes(), nil
	}
	if err := binary.Write(buf, binary.LittleEndian, []byte{'P','R','E','V'}); err != nil{
		return nil, err
	}
	/*
		binaryWriter.Write(previewSlot.HasData);
		binaryWriter.Write(previewSlot.Gil);
		binaryWriter.Write(previewSlot.PlayDuration);
		binaryWriter.Write(previewSlot.win_type);
	 */
	return buf.Bytes(), nil
}

