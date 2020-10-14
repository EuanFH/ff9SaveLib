package FF9Save

import (
	"bytes"
	"io/ioutil"
	"testing"
)

//need to properly test json unmarshaller with all possible save file locations
//probably just symlink the same file a bunch of times to save space

func TestSaveData_MarshalJsonFiles(t *testing.T) {

}

func TestSaveData_UnmarshalJsonFiles(t *testing.T) {
}

func TestSaveData_MarshalBinary(t *testing.T) {
	savedData := NewSavedData()
	if err := savedData.UnmarshalJsonFiles("testdata/JsonSavedData"); err != nil {
		t.Errorf("failed to unmarshal json files to savedData struct\n error: %s", err)
	}
	savedDataBytes, err := savedData.MarshalBinary()
	if err != nil {
		t.Errorf("failed to marshal binary from savedData struct\n error: %s", err)
	}
	savedDataOrigBytes, err := ioutil.ReadFile("testdata/SavedData_ww.dat")
	if err != nil {
		t.FailNow()
	}
	if !bytes.Equal(savedDataBytes, savedDataOrigBytes) {
		t.Error("binary generated from json does not match the original binary")
	}
}

func TestSaveData_UnmarshalBinary(t *testing.T) {

}
