package FF9Save

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestFilePreview_MarshalBinary(t *testing.T) {
	var filePreview FilePreview
	filePreviewJsonBytes, err := ioutil.ReadFile("testdata/filePreview.json")
	if err != nil {
		t.FailNow()
	}
	if err := json.Unmarshal(filePreviewJsonBytes, &filePreview); err != nil{
		t.FailNow()
	}
	filePreviewBinBytes, err := filePreview.MarshalBinary()
	if err != nil{
		t.FailNow()
	}
	if len(filePreviewBinBytes) != FilePreviewSize {
		t.Error("Binary generated from Json is too small")
	}
	/*
	if !bytes.Equal(fileBinBytes, fileOrigBinBytes) {
		t.Error("Binary generated from JSON does not match the original Binary")
	}
	 */
}

func TestFilePreview_Empty_MarshalBinary(t *testing.T) {
	var filePreview FilePreview
	filePreviewBinBytes, err := filePreview.MarshalBinary()
	if err != nil{
		t.FailNow()
	}
	if len(filePreviewBinBytes) != 965 {
		t.Error("Binary generated from Json is too small")
	}
	if !bytes.Equal(filePreviewBinBytes[:4], []byte{'N','O','N','E'}) {
		t.Error("File Preview does not start with NONE invalid empty save preview")
	}
}
