package main

import (
	"chinzer.net/ff9-save-converter/FF9Save"
	"io/ioutil"
	"log"
	"os"
)

func main(){
	args := os.Args[1:3]
	var switchSaveFolderPath string
	var binarySavePath string
	var convertToBinary bool
	for i, arg := range args {
		fileInfo, err := os.Stat(arg)
		if os.IsNotExist(err) {
			log.Fatalf("file %s is not found\n", arg)
		}
		if fileInfo.IsDir() {
			if i == 0 { convertToBinary = false}
			switchSaveFolderPath = arg
			continue
		}
		if i == 0 { convertToBinary = true}
		binarySavePath = arg
	}

	var saveData FF9Save.SaveData
	if convertToBinary {
		saveDataBytes, err := ioutil.ReadFile(binarySavePath)
		if err != nil{
			log.Fatalf("unable to read file %s\nerror: %s", binarySavePath, err)
		}
		if err := saveData.UnmarshalBinary(saveDataBytes); err != nil {
			log.Fatalf("failed to read file %s\nerror: %s", binarySavePath, err)
		}
		if err := saveData.MarshalJsonFiles(switchSaveFolderPath); err != nil {
			log.Fatalf("failed to generate json save files\nerror: %s", err)
		}
		return
	}
	if err := saveData.UnmarshalJsonFiles(switchSaveFolderPath); err != nil {
		log.Fatalf("failed to read json save files\nerror: %s", err)
	}
	saveDataBytes, err := saveData.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to create binary save file\n error: %s", err)
	}
	if err := ioutil.WriteFile(binarySavePath, saveDataBytes, 0644); err != nil {
		log.Fatalf("unable to write file to location %s\n error: %s", binarySavePath, err)
	}
}

