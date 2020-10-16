package FF9Save

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
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

func NewSavedData() SavedData {
	return SavedData{
		MetaData: NewMetaData(),
	}
}

func(sd *SavedData) MarshalJsonFiles(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
    	if err := os.Mkdir(directory, 0755); err != nil {
    		return err
		}
	}
	//File Info
	fileInfoBytes, err := json.Marshal(&sd.MetaData.FileInfo)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(directory + "/SLOTINFO", fileInfoBytes, 0644); err != nil {
		return err
	}

	var prefsLanguage PrefsLanguage
	languageString, ok := LanguageIntToLanguageString()[sd.MetaData.SelectedLanguage]
	if !ok {
		return err
	}
	prefsLanguage.Value = languageString
	prefsLanguageBytes, err := json.Marshal(prefsLanguage)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(directory +"/PREFS_Language", prefsLanguageBytes, 0644); err != nil {
		return err
	}

	//file previews
	for filePreviewPos, filePreview := range sd.FilePreviews {
		if filePreview == (FilePreview{}) {
			continue
		}
		fileName := generateSaveFileName("PREVIEW", filePreviewPos)
		filePreviewBytes, err := json.Marshal(&filePreview)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(directory + "/" + fileName, filePreviewBytes, 0644); err != nil {
			return err
		}
	}

	//auto save
	fileBytes, err := json.Marshal(&sd.Auto)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(directory + "/AUTO", fileBytes, 0644); err != nil {
		return err
	}

	//save files
	for filePos, file := range sd.Slot {
		if file == (File{}) {
			continue
		}
		fileName := generateSaveFileName("DATA", filePos)
		fileBytes, err := json.Marshal(&file)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(directory + "/" + fileName, fileBytes, 0644); err != nil {
			return err
		}
	}

	return nil
}

func(sd *SavedData) UnmarshalJsonFiles(directory string) error {
	//metadata
	fileInfoBytes, err := ioutil.ReadFile(directory +"/SLOTINFO")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(fileInfoBytes, &sd.MetaData.FileInfo); err != nil {
		return err
	}

	var prefsLanguage PrefsLanguage
	prefsLanguageBytes, err := ioutil.ReadFile(directory +"/PREFS_Language")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(prefsLanguageBytes, &prefsLanguage); err != nil {
		return err
	}
	languageInt, ok := LanguageStringToLanguageInt()[prefsLanguage.Value]
	if !ok {
		return fmt.Errorf("Language dosn't exist")
	}
	sd.MetaData.SelectedLanguage = languageInt

	//file previews
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, file := range files {
		fileNo := getFilePos("PREVIEW", file.Name())
		if fileNo == -1 || fileNo > len(sd.Slot) {
			continue
		}
		previewBytes, err := ioutil.ReadFile(directory + "/" + file.Name())
		if err != nil {
			return err
		}
		if err := json.Unmarshal(previewBytes, &sd.FilePreviews[fileNo]); err != nil {
			return err
		}

	}

	//auto save
	fileBytes, err := ioutil.ReadFile(directory + "/AUTO")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(fileBytes, &sd.Auto); err != nil {
		return err
	}

	//save files
	for _, file := range files {
		fileNo := getFilePos("DATA", file.Name())
		if fileNo == -1 {
			continue
		}
		fileBytes, err := ioutil.ReadFile(directory + "/" + file.Name())
		if err != nil {
			return err
		}
		if err := json.Unmarshal(fileBytes, &sd.Slot[fileNo]); err != nil {
			return err
		}

	}
	return nil
}

func getFilePos(prefix string, fileName string) int{
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

func getSlotAndFileNumber(filePos int) (int, int){
	slotNo := int(math.Floor(float64(filePos / FilesPerSlot)))
	fileNo := filePos - slotNo * FilesPerSlot
	return slotNo, fileNo
}

func generateSaveFileName(prefix string, filePos int) string{
	slotNo, fileNo := getSlotAndFileNumber(filePos)
	return fmt.Sprintf("%s_SLOT%d_FILE%d", prefix, slotNo, fileNo)
}

func(sd *SavedData) BinaryUnmarshaler(data []byte) error{
	buf := bytes.NewBuffer(data)
	//MetaData
	metaDataBytes, err := DecryptAndReadSaveSection(buf, MetaDataSize, MetaDataReservedSize)
	if err != nil{
		return err
	}
	if err := binary.Read(bytes.NewBuffer(metaDataBytes[4:]), binary.LittleEndian, &sd.MetaData); err != nil{
		return err
	}
	//FilePreviews
	for i, _ := range sd.FilePreviews {
		filePreviewBytes, err := DecryptAndReadSaveSection(buf, FilePreviewSize, FilePreviewReservedSize)
		if err != nil {
			return err
		}
		if err := sd.FilePreviews[i].UnBinaryMarshaler(filePreviewBytes); err != nil {
			return err
		}
	}
	//Auto
	fileBytes, err := DecryptAndReadSaveSection(buf, FileSize, FileReservedSize)
	if err != nil {
		return err
	}
	if err := sd.Auto.UnBinaryMarshaler(fileBytes); err != nil {
		return err
	}
	//Slot
	for i, _ := range sd.Slot {
		fileBytes, err := DecryptAndReadSaveSection(buf, FileSize, FileReservedSize)
		if err != nil {
			return err
		}
		if err := sd.Slot[i].UnBinaryMarshaler(fileBytes); err != nil {
			return err
		}
	}
	return nil
}

func(sd *SavedData) BinaryMarshaler() ([]byte, error){
	buf := new(bytes.Buffer)
	//MetaData
	metaDataBytes, err := sd.MetaData.BinaryMarshaler()
	if err != nil{
		return nil, err
	}
	if err := EncryptAndWriteSaveSection(buf, metaDataBytes, MetaDataReservedSize); err != nil {
		return nil, err
	}
	//FilePreviews
	for _, filePreview := range sd.FilePreviews {
		filePreviewBytes, err := filePreview.BinaryMarshaler()
		if err != nil{
			return nil, err
		}
		if err := EncryptAndWriteSaveSection(buf, filePreviewBytes, FilePreviewReservedSize); err != nil {
			return nil, err
		}
	}
	//Auto
	fileBytes, err := sd.Auto.BinaryMarshaler()
	if err != nil{
		return nil, err
	}
	if err := EncryptAndWriteSaveSection(buf, fileBytes, FileReservedSize); err != nil {
		return nil, err
	}
	//Slot
	for i, _ := range sd.Slot {
		fileBytes, err := sd.Slot[i].BinaryMarshaler()
		if err != nil{
			return nil, err
		}
		if err := EncryptAndWriteSaveSection(buf, fileBytes, FileReservedSize); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func EncryptAndWriteSaveSection(buf *bytes.Buffer, data []byte, reservedSize int) error{
	dataEncrypted, err := Encrypt(data)
	if err != nil{
		return err
	}
	dataEncryptedWithReservedSpace := append(dataEncrypted, make([]byte, reservedSize - len(dataEncrypted))...)
	if err := binary.Write(buf, binary.LittleEndian, dataEncryptedWithReservedSpace); err != nil {
		return err
	}
	return nil
}

func DecryptAndReadSaveSection(buf *bytes.Buffer, size, reservedSize int) ([]byte, error){
	dataEncrypted := make([]byte, cipherSize(size))
	if _, err := buf.Read(dataEncrypted); err != nil{
		return nil, err
	}
	data, err := Decrypt(dataEncrypted)
	if err != nil{
		return nil, err
	}
	buf.Next(reservedSize - len(dataEncrypted))
	return data, nil
}
