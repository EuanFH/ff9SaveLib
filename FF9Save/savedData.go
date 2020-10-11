package FF9Save

import (
	"bytes"
	"chinzer.net/ff9-save-converter/FF9Save/Crypto"
	"encoding/binary"
)

const numSlots=10
const numFiles=15

type SaveData struct {
	MetaData MetaData
	FilePreviews [150]FilePreview
	Auto File
	Slot [150]File
}

func NewSaveData() SaveData{
	return SaveData{
		MetaData: NewMetaData(),
	}
}

func(sd *SaveData) UnmarshalBinary(data []byte) error{
	buf := bytes.NewBuffer(data)
	//MetaData
	metaDataBytes, err := Crypto.DecryptAndReadSaveSection(buf, MetaDataSize, MetaDataReservedSize)
	if err != nil{
		return err
	}
	if err := binary.Read(bytes.NewBuffer(metaDataBytes[4:]), binary.LittleEndian, &sd.MetaData); err != nil{
		return err
	}
	//FilePreviews
	for i, _ := range sd.FilePreviews {
		filePreviewBytes, err := Crypto.DecryptAndReadSaveSection(buf, FilePreviewSize, FilePreviewReservedSize)
		if err != nil {
			return err
		}
		if err := sd.FilePreviews[i].UnmarshalBinary(filePreviewBytes); err != nil {
			return err
		}
	}
	//Auto
	fileBytes, err := Crypto.DecryptAndReadSaveSection(buf, FileSize, FileReservedSize)
	if err != nil {
		return err
	}
	if err := sd.Auto.UnmarshalBinary(fileBytes); err != nil {
		return err
	}
	//Slot
	for i, _ := range sd.Slot {
		fileBytes, err := Crypto.DecryptAndReadSaveSection(buf, FileSize, FileReservedSize)
		if err != nil {
			return err
		}
		if err := sd.Slot[i].UnmarshalBinary(fileBytes); err != nil {
			return err
		}
	}
	return nil
}

func(sd *SaveData) MarshalBinary() ([]byte, error){
	buf := new(bytes.Buffer)
	//MetaData
	metaDataBytes, err := sd.MetaData.MarshalBinary()
	if err != nil{
		return nil, err
	}
	if err := Crypto.EncryptAndWriteSaveSection(buf, metaDataBytes, MetaDataReservedSize); err != nil {
		return nil, err
	}
	//FilePreviews
	for _, filePreview := range sd.FilePreviews {
		filePreviewBytes, err := filePreview.MarshalBinary()
		if err != nil{
			return nil, err
		}
		if err := Crypto.EncryptAndWriteSaveSection(buf, filePreviewBytes, FilePreviewReservedSize); err != nil {
			return nil, err
		}
	}
	//Auto
	fileBytes, err := sd.Auto.MarshalBinary()
	if err != nil{
		return nil, err
	}
	if err := Crypto.EncryptAndWriteSaveSection(buf, fileBytes, FileReservedSize); err != nil {
		return nil, err
	}
	//Slot
	for i, _ := range sd.Slot {
		fileBytes, err := sd.Slot[i].MarshalBinary()
		if err != nil{
			return nil, err
		}
		if err := Crypto.EncryptAndWriteSaveSection(buf, fileBytes, FileReservedSize); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}


