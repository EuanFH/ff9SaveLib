package FF9Save

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/euanfh/ff9SaveLib/Crypto"
	"io/ioutil"
	"regexp"
	"strconv"
)

const MaxSlots=10
const FilesPerSlot=15

type SavedData struct {
	MetaData MetaData
	FilePreviews [150]FilePreview
	Auto File
	Slot [150]File
}

//if i put into the unmarshal function to set these fields up properly i wouldnt need this
func NewSavedData() SavedData {
	return SavedData{
		MetaData: NewMetaData(),
	}
}

func(sd *SavedData) MarshalJsonFiles(directory string) error {
	return nil
}

func(sd *SavedData) UnmarshalJsonFiles(directory string) error {
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
		fileNo := getFileNumber("PREVIEW", file.Name())
		if fileNo == -1 {
			continue
		}
		previewBytes, err := ioutil.ReadFile(directory + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(previewBytes, &sd.FilePreviews[fileNo]); err != nil {
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
		fileNo := getFileNumber("DATA", file.Name())
		if fileNo == -1 {
			continue
		}
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

func getFileNumber(prefix string, fileName string) int{
	re := regexp.MustCompile(fmt.Sprintf(`%s_SLOT(?P<SLOT>\d{1,2})_FILE(?P<FILE>\d{1,2})$`, prefix))
	slotFile := re.FindStringSubmatch(fileName)
	if slotFile == nil {
		return -1
	}
	if len(slotFile) < 3 {
		return -1
	}
	slotFile = slotFile[1:]
	slotNo, err := strconv.Atoi(slotFile[0])
	if err != nil {
		return -1
	}
	fileNo, err := strconv.Atoi(slotFile[1])
	if err != nil {
		return -1
	}
	return slotNo * FilesPerSlot + fileNo
}

func(sd *SavedData) UnmarshalBinary(data []byte) error{
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

func(sd *SavedData) MarshalBinary() ([]byte, error){
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


