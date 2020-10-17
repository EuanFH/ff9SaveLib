package main

import (
        "encoding/json"
        "fmt"
        "github.com/euanfh/ff9SaveLib"
        "io/ioutil"
        "os"
        "path/filepath"
)

func main(){
        if len(os.Args) < 3 {
                helpMessage()
                return
        }
        args := os.Args[1:3]
        if args[0] == "-h" || args[0] == "--help" {
                helpMessage()
                return
        }

        var switchSaveFolderPath string
        var binarySavePath string

        firstFileInfo, err := os.Stat(args[0])
        if os.IsNotExist(err) {
                fmt.Printf("file %s is not found\n", args[0])
                os.Exit(1)
        }
        if firstFileInfo.IsDir() {
                switchSaveFolderPath = filepath.Clean(args[0])
                binarySavePath = filepath.Clean(args[1])
                if err := convertJsonSavesToBinary(switchSaveFolderPath, binarySavePath); err != nil {
                        fmt.Print(err)
                        os.Exit(1)
                }
                return
        }
        binarySavePath = filepath.Clean(args[0])
        switchSaveFolderPath = filepath.Clean(args[1])
        switchSaveFolderInfo, err := os.Stat(switchSaveFolderPath)
        if os.IsNotExist(err) {
                if err := os.MkdirAll(switchSaveFolderPath, 0755); err != nil {
                        fmt.Printf("failed to make directory %s", switchSaveFolderPath)
                        os.Exit(1)
                }

                switchSaveFolderInfo, err = os.Stat(switchSaveFolderPath)
                if err != nil {
                        fmt.Printf("failed to make directory %s", switchSaveFolderPath)
                        os.Exit(1)
                }
        }
        if !switchSaveFolderInfo.IsDir() {
                fmt.Printf("cannot write json save files path provided is not a directory")
                os.Exit(1)
        }
        if err := convertBinarySaveToJson(binarySavePath, switchSaveFolderPath); err != nil{
                fmt.Print(err)
                os.Exit(1)
        }
}

func helpMessage() {
        fmt.Printf("Final Fantasy 9 Save Converter\n" +
                "Usage:\n" +
                "\tff9sc FILE DIRECTORY\n" +
                "\tff9sc DIRECTORY FILE\n" +
                "Examples:\n" +
                "\t$ ff9sc SavedData_ww.dat switchSavesFolder\n" +
                "\t$ ff9sc switchSavesFolder SavedData_ww.dat\n")
}

func convertBinarySaveToJson(binarySavePath string, switchSaveFolderPath string) error{
        savedData := ff9Save.NewSavedData()
        saveDataBytes, err := ioutil.ReadFile(binarySavePath)
        if err != nil{
                return fmt.Errorf("unable to read file %s\nerror: %s", binarySavePath, err)
        }
        if err := savedData.BinaryUnmarshaler(saveDataBytes); err != nil {
                return fmt.Errorf("failed to read file %s\nerror: %s", binarySavePath, err)
        }
        if err := marshalFF9JsonFiles(savedData, switchSaveFolderPath); err != nil {
                return fmt.Errorf("failed to generate json save files\nerror: %s", err)
        }
        return nil
}
func convertJsonSavesToBinary(switchSaveFolderPath string, binarySavePath string) error {
        savedData, err := unmarshalFF9JsonFiles(switchSaveFolderPath)
        if err != nil {
                return fmt.Errorf("failed to read json save files\nerror: %s", err)
        }
        saveDataBytes, err := savedData.BinaryMarshaler()
        if err != nil {
                return fmt.Errorf("failed to create binary save file\n error: %s", err)
        }
        if err := ioutil.WriteFile(binarySavePath, saveDataBytes, 0644); err != nil {
                return fmt.Errorf("unable to write file to location %s\n error: %s", binarySavePath, err)
        }
        return nil
}


func marshalFF9JsonFiles(savedData ff9Save.SavedData, directory string) error {
        //File Info
        fileInfoBytes, err := json.Marshal(&savedData.MetaData.FileInfo)
        if err != nil {
                return err
        }

        if err := ioutil.WriteFile(filepath.Join(directory, "SLOTINFO"), fileInfoBytes, 0644); err != nil {
                return err
        }

        var prefsLanguage ff9Save.PrefsLanguage
        languageString, ok := ff9Save.LanguageIntToLanguageString()[savedData.MetaData.SelectedLanguage]
        if !ok {
                return err
        }
        prefsLanguage.Value = languageString
        prefsLanguageBytes, err := json.Marshal(prefsLanguage)
        if err != nil {
                return err
        }
        if err := ioutil.WriteFile(filepath.Join(directory, "PREFS_Language"), prefsLanguageBytes, 0644); err != nil {
                return err
        }

        //file previews
        for filePreviewPos, filePreview := range savedData.FilePreviews {
                if filePreview == (ff9Save.FilePreview{}) {
                        continue
                }
                fileName := ff9Save.GenerateSaveFileName("PREVIEW", filePreviewPos)
                filePreviewBytes, err := json.Marshal(&filePreview)
                if err != nil {
                        return err
                }
                if err := ioutil.WriteFile(filepath.Join(directory, fileName), filePreviewBytes, 0644); err != nil {
                        return err
                }
        }

        //auto save
        fileBytes, err := json.Marshal(&savedData.Auto)
        if err != nil {
                return err
        }

        if err := ioutil.WriteFile(filepath.Join(directory, "AUTO"), fileBytes, 0644); err != nil {
                return err
        }

        //save files
        for filePos, file := range savedData.Slot {
                if file == (ff9Save.File{}) {
                        continue
                }
                fileName := ff9Save.GenerateSaveFileName("DATA", filePos)
                fileBytes, err := json.Marshal(&file)
                if err != nil {
                        return err
                }
                if err := ioutil.WriteFile(filepath.Join(directory, fileName), fileBytes, 0644); err != nil {
                        return err
                }
        }

        return nil
}

func unmarshalFF9JsonFiles(directory string) (*ff9Save.SavedData, error) {
        savedData := ff9Save.NewSavedData()
        //metadata
        fileInfoBytes, err := ioutil.ReadFile(filepath.Join(directory, "SLOTINFO"))
        if err != nil {
                return nil, err
        }
        if err := json.Unmarshal(fileInfoBytes, &savedData.MetaData.FileInfo); err != nil {
                return nil, err
        }

        var prefsLanguage ff9Save.PrefsLanguage
        prefsLanguageBytes, err := ioutil.ReadFile(filepath.Join(directory, "/PREFS_Language"))
        if err != nil {
                return nil, err
        }
        if err := json.Unmarshal(prefsLanguageBytes, &prefsLanguage); err != nil {
                return nil, err
        }
        languageInt, ok := ff9Save.LanguageStringToLanguageInt()[prefsLanguage.Value]
        if !ok {
                return nil, fmt.Errorf("language dosn't exist")
        }
        savedData.MetaData.SelectedLanguage = languageInt

        //file previews
        files, err := ioutil.ReadDir(directory)
        if err != nil {
                return nil, err
        }
        for _, file := range files {
                fileNo := ff9Save.GetFilePos("PREVIEW", file.Name())
                if fileNo == -1 || fileNo > len(savedData.Slot) {
                        continue
                }
                previewBytes, err := ioutil.ReadFile(filepath.Join(directory, file.Name()))
                if err != nil {
                        return nil, err
                }
                if err := json.Unmarshal(previewBytes, &savedData.FilePreviews[fileNo]); err != nil {
                        return nil, err
                }

        }

        //auto save
        fileBytes, err := ioutil.ReadFile(filepath.Join(directory, "AUTO"))
        if err != nil {
                return nil, err
        }
        if err := json.Unmarshal(fileBytes, &savedData.Auto); err != nil {
                return nil, err
        }

        //save files
        for _, file := range files {
                fileNo := ff9Save.GetFilePos("DATA", file.Name())
                if fileNo == -1 {
                        continue
                }
                fileBytes, err := ioutil.ReadFile(filepath.Join(directory, file.Name()))
                if err != nil {
                        return nil, err
                }
                if err := json.Unmarshal(fileBytes, &savedData.Slot[fileNo]); err != nil {
                        return nil, err
                }

        }
        return &savedData, nil
}