package FF9Save

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

type String128 [128]byte

func (s String128) MarshalJSON() ([]byte, error) {
	strBytes := []byte(string(s[:]))
	//remove null bytes
	strBytes = bytes.Trim(strBytes, "\x00")
	quotedStrBytes := []byte(`"` + string(strBytes) + `"`)
	return quotedStrBytes, nil
}

func(s *String128) UnmarshalJSON(data []byte) error{
	var str string
	if err := json.Unmarshal(data, &str); err != nil{
		return err
	}
	for i, c := range []byte(str) {
		s[i] = c
	}
	return nil
}

type String4K [4096]byte

func (s String4K) MarshalJSON() ([]byte, error) {
	strBytes := []byte(string(s[:]))
	//remove null bytes
	strBytes = bytes.Trim(strBytes, "\x00")
	quotedStrBytes := []byte(`"` + string(strBytes) + `"`)
	return quotedStrBytes, nil
}

func(s *String4K) UnmarshalJSON(data []byte) error{
	var str string
	if err := json.Unmarshal(data, &str); err != nil{
		return err
	}
	for i, c := range []byte(str) {
		s[i] = c
	}
	return nil
}

type byteFromStr byte

func (b byteFromStr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(b), 10) + `"`), nil
}

func (b *byteFromStr) UnmarshalJSON(data []byte) error {
	// Try array of strings first.
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(str, 10, 8)
	if err != nil{
		return err
	}
	*b = byteFromStr(i)
	return nil
}

type uint32FromStr uint32

func (ui uint32FromStr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(ui), 10) + `"`), nil
}

func (ui *uint32FromStr) UnmarshalJSON(data []byte) error {
	// Try array of strings first.
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil{
		return err
	}
	*ui = uint32FromStr(i)
	return nil
}

type int32FromStr int32

func (in int32FromStr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(in), 10) + `"`), nil
}

func (in *int32FromStr) UnmarshalJSON(data []byte) error {
	// Try array of strings first.
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return err
	}
	*in = int32FromStr(i)
	return nil
}

type uint16FromStr uint16

func (ui uint16FromStr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(ui), 10) + `"`), nil
}

func (ui *uint16FromStr) UnmarshalJSON(data []byte) error {
	// Try array of strings first.
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil{
		return err
	}
	*ui = uint16FromStr(i)
	return nil
}

type boolFromStr bool

func (b boolFromStr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strings.Title(strconv.FormatBool(bool(b))) + `"`), nil
}

func (b *boolFromStr) UnmarshalJSON(data []byte) error {
	// Try array of strings first.
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	bParsed, err := strconv.ParseBool(str)
	if err != nil{
		return err
	}
	*b = boolFromStr(bParsed)
	return nil
}
