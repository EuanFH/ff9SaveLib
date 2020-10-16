package FF9Save

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"github.com/d1str0/pkcs7"
	"golang.org/x/crypto/pbkdf2"
)

var password = []byte("System.Security.SecureString")
var salt = []byte{3,3,1,4,7,0,9,7}
const AES256KeySize = 32

//Encrypt with AES CBC mode using the ff9 password and salt including pkcs7 padding
func Encrypt(data []byte) ([]byte, error) {
	data, err := pkcs7.Pad(data, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	key, iv := getKeyAndIV()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	cipheredData := make([]byte, len(data))
	mode.CryptBlocks(cipheredData, data)
	return cipheredData, nil
}

//Decrypt with AES CBC mode using the ff9 password and salt removing pkcs7 padding
func Decrypt(data []byte) ([]byte, error){
	key, iv := getKeyAndIV()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	cipheredData := make([]byte, len(data))
	mode.CryptBlocks(cipheredData, data)
	cipheredData, err = pkcs7.Unpad(cipheredData)
	if err != nil {
		return nil, err
	}
	return cipheredData, nil
}

//Key is generated the normal way. IV is the 16 bytes after key generation
func getKeyAndIV() ([]byte, []byte){
	keyIV := pbkdf2.Key(password, salt, 1000, AES256KeySize + aes.BlockSize, sha1.New)
	return keyIV[:AES256KeySize], keyIV[AES256KeySize:]
}

//data size to the nearest AES block
func cipherSize(size int) int {
	const blocksize = 16
	return size + blocksize - size % blocksize
}
