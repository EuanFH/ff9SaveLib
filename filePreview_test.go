package ff9Save

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestFilePreviewBinaryMarshaler(t *testing.T) {
	var filePreview FilePreview
	filePreviewJsonBytes, err := ioutil.ReadFile("testdata/filePreview.json")
	if err != nil {
		t.FailNow()
	}
	if err := json.Unmarshal(filePreviewJsonBytes, &filePreview); err != nil{
		t.FailNow()
	}
	filePreviewBinBytes, err := filePreview.BinaryMarshaler()
	if err != nil{
		t.FailNow()
	}
	if len(filePreviewBinBytes) != FilePreviewSize {
		t.Error("Binary generated from Json is too small")
	}

	/*if !bytes.Equal(fileBinBytes, fileOrigBinBytes) {
		t.Error("Binary generated from JSON does not match the original Binary")
	}*/
}

func TestFilePreviewEmptyBinaryMarshaler(t *testing.T) {
	var filePreview FilePreview
	filePreviewBinBytes, err := filePreview.BinaryMarshaler()
	if err != nil{
		t.FailNow()
	}
	if len(filePreviewBinBytes) != FilePreviewSize {
		t.Error("Binary generated from Json is too small")
	}
	if !bytes.Equal(filePreviewBinBytes[:4], NoneHeader[:]) {
		t.Error("File Preview does not start with NONE invalid empty save preview")
	}
}

func TestFilePreviewFixCharacterInfoFromBinary(t *testing.T) {
	var filePreview FilePreview
	filePreview.CharacterInfoList[0].SerialID = -1
	filePreview.CharacterInfoList[0].Level = -1
	filePreview.FixCharacterInfoFromBinary()
	if filePreview.CharacterInfoList[0].SerialID != 0 || filePreview.CharacterInfoList[0].Level != 0 {
		t.Error("Empty character info incorrect converted from binary")
	}
}

func TestFilePreview_FixCharacterInfoForBinary(t *testing.T) {
	var filePreview FilePreview
	filePreview.FixCharacterInfoForBinary()
	if filePreview.CharacterInfoList[0].SerialID != -1 || filePreview.CharacterInfoList[0].Level != -1 {
		t.Error("Empty character info incorrect converted from binary")
	}
}
