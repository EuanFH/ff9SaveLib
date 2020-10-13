package FF9Save

import (
	"bytes"
	"chinzer.net/ff9-save-converter/FF9Save/Crypto"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
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

func(sd *SaveData) MarshalJsonFiles(directory string) error {
	return nil
}

func(sd *SaveData) UnmarshalJsonFiles(directory string) error {
	fileInfoBytes, err := ioutil.ReadFile(directory +"/SLOTINFO")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(fileInfoBytes, &sd.MetaData.FileInfo); err != nil {
		panic(err)
	}
	sd.MetaData.SelectedLanguage = 1 //cba reading file this now
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if(!strings.HasPrefix(file.Name(), "PREVIEW")){
			continue
		}
		slotRune := []rune(file.Name())[12]
		slotNo, err := strconv.Atoi(string(slotRune))
		if err != nil{
			panic(err)
		}
		fileRune := []rune(file.Name())[18]
		fileNo, err := strconv.Atoi(string(fileRune))
		if err != nil{
			panic(err)
		}

		previewFileNo := slotNo * 15 + fileNo


		previewBytes, err := ioutil.ReadFile(directory + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(previewBytes, &sd.FilePreviews[previewFileNo]); err != nil {
			panic(err)
		}

	}
	fileBytes, err := ioutil.ReadFile(directory + "/AUTO")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(fileBytes, &sd.Auto); err != nil {
		panic(err)
	}
	for _, file := range files {
		if(!strings.HasPrefix(file.Name(), "DATA")){
			continue
		}
		slotRune := []rune(file.Name())[9]
		slotNo, err := strconv.Atoi(string(slotRune))
		if err != nil{
			panic(err)
		}
		fileRune := []rune(file.Name())[15]
		fileNo, err := strconv.Atoi(string(fileRune))
		if err != nil{
			panic(err)
		}

		fileNo = slotNo * 15 + fileNo


		fileBytes, err := ioutil.ReadFile(directory + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(fileBytes, &sd.Slot[fileNo]); err != nil {
			panic(err)
		}

	}
	return nil
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


