package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha1"
    "golang.org/x/crypto/pbkdf2"
)

//assuming this is extracted from the game
//hex decode?
var password = []byte("67434cd0-1ca3-11e5-9a21-1697f925ec7b7a5313a0-1ca3-11e5-b939-0800200c9a66")
var salt = []byte{3,3,1,4,7,0,9,7}
const AES256KeySize = 32

const saveBlockSize = 18432

func main(){
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
