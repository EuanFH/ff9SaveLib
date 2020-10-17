package ff9Save

import (
	"bytes"
	"testing"
)

func TestMetaData_BinaryMarshaler(t *testing.T) {
	metaData := NewMetaData()
	metaDataBytes, err := metaData.BinaryMarshaler()
	if err != nil{
		t.FailNow()
	}
	if(!bytes.Equal(metaDataBytes[:4], SaveHeader[:])){
		t.Errorf("metaData binary contains incorrect header")
	}
	if len(metaDataBytes) != MetaDataSize {
		t.Error("metaData binary generated is too small")
	}
}
