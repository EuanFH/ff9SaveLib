package FF9Save

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestFile_BinaryMarshaler(t *testing.T) {
	var file File
	fileJsonBytes, err := ioutil.ReadFile("testdata/file.json")
	if err != nil {
		t.FailNow()
	}
	if err := json.Unmarshal(fileJsonBytes, &file); err != nil{
		t.FailNow()
	}
	fileBinBytes, err := file.BinaryMarshaler()
	if err != nil{
		t.FailNow()
	}
	fileOrigBinBytes, err := ioutil.ReadFile("testdata/file.bin")
	if err != nil {
		t.FailNow()
	}
	//test the size
	if !bytes.Equal(fileBinBytes, fileOrigBinBytes) {
		t.Error("Binary generated from JSON does not match the original Binary")
	}
}

func TestFile_UnBinaryMarshaler(t *testing.T) {
	var file File
	fileBinBytes, err := ioutil.ReadFile("testdata/file.bin")
	if err != nil {
		t.FailNow()
	}
	if err := file.UnBinaryMarshaler(fileBinBytes); err != nil{
		t.Fail()
	}
	fileJsonBytes, err := json.Marshal(file)
	if err != nil {
		t.FailNow()
	}
	fileOrigJsonBytes, err := ioutil.ReadFile("testdata/file.json")
	if err != nil{
		t.FailNow()
	}
	if !bytes.Equal(fileJsonBytes, fileOrigJsonBytes) {
		t.Error("Binary generated from JSON does not match the original Binary")
	}
}

