package main

import (
	"chinzer.net/ff9-save-converter/FF9Save"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

func main(){
	saveData := readSwitchSave()
	saveDataBinary, err := saveData.MarshalBinary()
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("saveDataFromSwitch", saveDataBinary, 0777); err != nil {
		panic(err)
	}
	var saveDataOrig FF9Save.SaveData
	saveDataBytes, err := ioutil.ReadFile("SavedData_ww.dat")
	if err != nil{
		panic(err)
	}
	if err := saveDataOrig.UnmarshalBinary(saveDataBytes); err != nil {
		panic(err)
	}
	/*
	var saveData FF9Save.SaveData
	saveDataBytes, err := ioutil.ReadFile("SavedData_ww.dat")
	if err != nil{
		panic(err)
	}
	if err := saveData.UnmarshalBinary(saveDataBytes); err != nil {
		panic(err)
	}
	saveDataNew, err := saveData.MarshalBinary()
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("newSaveData", saveDataNew, 0777); err != nil {
		panic(err)
	}
	 */
    /*
    jsonFile, err := os.Open("origSave.json")
    // if we os.Open returns an error then handle it
    if err != nil {
    	panic(err)
    }
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()
	ff9JsonSaveBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
	    panic(err)
    }
    var ff9Save FF9Save.Slot
    if err := json.Unmarshal(ff9JsonSaveBytes, &ff9Save); err != nil {
        panic(err)
    }

    ff9JsonString, err := json.Marshal(ff9Save)
    if err != nil {
        panic(err)
    }
    err = ioutil.WriteFile("output.json", ff9JsonString, 0777)
    if err != nil{
        panic(err)
    }

    ff9SaveBytes, err := ff9Save.MarshalBinary()
    if err != nil {
        panic(err)
    }

    err = ioutil.WriteFile("output.bin", ff9SaveBytes, 0777)
    if err != nil{
        panic(err)
    }
     */
}

//this code will not work in a lot of cases
func readSwitchSave() FF9Save.SaveData{
	saveData := FF9Save.NewSaveData()
	fileInfoBytes, err := ioutil.ReadFile("SwitchSaves/SLOTINFO")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(fileInfoBytes, &saveData.MetaData.FileInfo); err != nil {
		panic(err)
	}
	saveData.MetaData.SelectedLanguage = 1 //cba reading file this now
	files, err := ioutil.ReadDir("SwitchSaves")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if(!strings.HasPrefix(file.Name(), "PREVIEW")){
			continue
		}
		slotRune := []rune(file.Name())[12]
		slotNo, err := strconv.Atoi(string(slotRune))
		if err != nil{
			panic(err)
		}
		fileRune := []rune(file.Name())[18]
		fileNo, err := strconv.Atoi(string(fileRune))
		if err != nil{
			panic(err)
		}

		previewFileNo := slotNo * 15 + fileNo


		previewBytes, err := ioutil.ReadFile("SwitchSaves/" + file.Name())
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(previewBytes, &saveData.FilePreviews[previewFileNo]); err != nil {
			panic(err)
		}

	}
	fileBytes, err := ioutil.ReadFile("SwitchSaves/AUTO")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(fileBytes, &saveData.Auto); err != nil {
		panic(err)
	}
	for _, file := range files {
		if(!strings.HasPrefix(file.Name(), "DATA")){
			continue
		}
		slotRune := []rune(file.Name())[9]
		slotNo, err := strconv.Atoi(string(slotRune))
		if err != nil{
			panic(err)
		}
		fileRune := []rune(file.Name())[15]
		fileNo, err := strconv.Atoi(string(fileRune))
		if err != nil{
			panic(err)
		}

		fileNo = slotNo * 15 + fileNo


		fileBytes, err := ioutil.ReadFile("SwitchSaves/" + file.Name())
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(fileBytes, &saveData.Slot[fileNo]); err != nil {
			panic(err)
		}

	}
	return saveData
}
