package FF9Save

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestSlot_MarshalBinary(t *testing.T) {
	var slot Slot
	slotJsonBytes, err := ioutil.ReadFile("testdata/slot.json")
	if err != nil {
		t.FailNow()
	}
	if err := json.Unmarshal(slotJsonBytes, &slot); err != nil{
		t.FailNow()
	}
	slotBinBytes, err := slot.MarshalBinary()
	if err != nil{
		t.FailNow()
	}
	slotOrigBinBytes, err := ioutil.ReadFile("testdata/slot.bin")
	if err != nil {
		t.FailNow()
	}
	//test the size
	if !bytes.Equal(slotBinBytes, slotOrigBinBytes) {
		t.Error("Binary generated from JSON does not match the original Binary")
	}
}

func TestSlot_UnmarshalBinary(t *testing.T) {
	var slot Slot
	slotBinBytes, err := ioutil.ReadFile("testdata/slot.bin")
	if err != nil {
		t.FailNow()
	}
	if err := slot.UnmarshalBinary(slotBinBytes); err != nil{
		t.Fail()
	}
	slotJsonBytes, err := json.Marshal(slot)
	if err != nil {
		t.FailNow()
	}
	slotOrigJsonBytes, err := ioutil.ReadFile("testdata/slot.json")
	if err != nil{
		t.FailNow()
	}
	if !bytes.Equal(slotJsonBytes, slotOrigJsonBytes) {
		t.Error("Binary generated from JSON does not match the original Binary")
	}
}

