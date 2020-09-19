package main

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha1"
    "encoding/binary"
    "encoding/json"
    "golang.org/x/crypto/pbkdf2"
    "io/ioutil"
    "os"
)

//assuming this is extracted from the game
//hex decode?
var password = []byte("67434cd0-1ca3-11e5-9a21-1697f925ec7b7a5313a0-1ca3-11e5-b939-0800200c9a66")
var salt = []byte{3,3,1,4,7,0,9,7}
const AES256KeySize = 32

const saveBlockSize = 18432

func main(){
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
    var ff9Save FF9Save
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

    //blyt to file
    buf := new(bytes.Buffer)

    if err := binary.Write(buf, binary.LittleEndian, []byte{'S','A','V','E'}); err != nil{
        panic(err)
    }
    if err := binary.Write(buf, binary.LittleEndian, ff9Save); err != nil{
        panic(err)
    }
    err = ioutil.WriteFile("output.bin", buf.Bytes(), 0777)
    if err != nil{
        panic(err)
    }

    //read and put into json
    var ff9SaveFromBin FF9Save
    binFile, err := os.Open("origSave.bin")
    if err != nil {
        panic(err)
    }
    defer binFile.Close()
    ff9BinSaveBytes, err := ioutil.ReadAll(binFile)
    if err != nil {
        panic(err)
    }

    ff9BinSaveBytes = ff9BinSaveBytes[4:]
    ff9BinSaveBytes = append(ff9BinSaveBytes, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
    binBuf := bytes.NewReader(ff9BinSaveBytes)
    if err := binary.Read(binBuf, binary.LittleEndian, &ff9BinSaveBytes); err != nil {
    	panic(err)
    }

    newJsonSaveBuf, err := json.Marshal(ff9SaveFromBin)
    if err := ioutil.WriteFile("binToJson.json", newJsonSaveBuf, os.ModePerm); err != nil{
        panic(err)
    }
    /*
  //filePath := os.Args[1]
    filePath := "/home/chinz/mnt/HentaiGames/SteamGames/steamapps/compatdata/377840/pfx/drive_c/users/steamuser/AppData/LocalLow/SquareEnix/FINAL FANTASY IX/Steam/EncryptedSavedData/SavedData_ww.dat"
    encryptedSaveFile, err := ioutil.ReadFile(filePath)
  if err != nil {
   fmt.Println("failed to read save file")
   return
  }
  decryptedSaveFile, err := decryptSaveFile(encryptedSaveFile)
  if err != nil{
    fmt.Println("failed to decrypt save file")
    return
  }
  if err = ioutil.WriteFile("./decryptedSave", decryptedSaveFile, 0777); err != nil{
    fmt.Println("failed to write decrypted save file")
    return
  }
    */
}

func decryptSaveFile(encryptedSave []byte) ([]byte, error){
  //don't know why its 1000 iterations
  key := pbkdf2.Key(password, salt, 1000, AES256KeySize, sha1.New)
  block, err := aes.NewCipher(key)
  if err != nil {
    return nil, err
  }
  mode := cipher.NewCBCDecrypter(block, key[:aes.BlockSize])
  decryptedSave := make([]byte, len(encryptedSave))
  mode.CryptBlocks(decryptedSave, encryptedSave)
  return decryptedSave, nil
}
