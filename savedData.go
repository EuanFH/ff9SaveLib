package ff9Save

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
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

func GetFilePos(prefix string, fileName string) int{
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

func GetSlotAndFileNumber(filePos int) (int, int){
	slotNo := int(math.Floor(float64(filePos / FilesPerSlot)))
	fileNo := filePos - slotNo * FilesPerSlot
	return slotNo, fileNo
}

func GenerateSaveFileName(prefix string, filePos int) string{
	slotNo, fileNo := GetSlotAndFileNumber(filePos)
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
	for i := range sd.FilePreviews {
		filePreviewBytes, err := DecryptAndReadSaveSection(buf, FilePreviewSize, FilePreviewReservedSize)
		if err != nil {
			return err
		}
		if err := sd.FilePreviews[i].BinaryUnmarshaler(filePreviewBytes); err != nil {
			return err
		}
	}
	//Auto
	fileBytes, err := DecryptAndReadSaveSection(buf, FileSize, FileReservedSize)
	if err != nil {
		return err
	}
	if err := sd.Auto.BinaryUnmarshaler(fileBytes); err != nil {
		return err
	}
	//Slot
	for i := range sd.Slot {
		fileBytes, err := DecryptAndReadSaveSection(buf, FileSize, FileReservedSize)
		if err != nil {
			return err
		}
		if err := sd.Slot[i].BinaryUnmarshaler(fileBytes); err != nil {
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
	for i := range sd.Slot {
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