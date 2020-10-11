package Crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"golang.org/x/crypto/pbkdf2"
)

var password = []byte("System.Security.SecureString")
var salt = []byte{3,3,1,4,7,0,9,7}
const AES256KeySize = 32

func EncryptAndWriteSaveSection(buf *bytes.Buffer, data []byte, reservedSize int) error{
	data, err := pkcs7Pad(data, aes.BlockSize)
	dataEncrypted, err := Encrypt(data)
	if err != nil{
		return err
	}
	dataEncryptedWithPadding := append(dataEncrypted, make([]byte, reservedSize - len(data))...)
	if err := binary.Write(buf, binary.LittleEndian, dataEncryptedWithPadding); err != nil {
		return err
	}
	return nil
	/*
	dataEncryptedWithPadding := append(data, make([]byte, reservedSize - len(data))...)
	if err := binary.Write(buf, binary.LittleEndian, dataEncryptedWithPadding); err != nil {
		return err
	}
	return nil
	*/
}

func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errors.New("yee")
	}
	if b == nil || len(b) == 0 {
		return nil, errors.New("yee")
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

func DecryptAndReadSaveSection(buf *bytes.Buffer, size, reservedSize int) ([]byte, error){
	dataEncrypted := make([]byte, CipherSize(size))
	if _, err := buf.Read(dataEncrypted); err != nil{
		return nil, err
	}
	data, err := Decrypt(dataEncrypted)
	if err != nil{
		return nil, err
	}
	buf.Next(reservedSize - len(dataEncrypted))
	return data, nil
}

func Encrypt(data []byte) ([]byte, error) {
	return EncryptDecrypt(data, true)
}

func Decrypt(encryptedData []byte) ([]byte, error){
	return EncryptDecrypt(encryptedData, false)
}

func EncryptDecrypt(data []byte, encrypt bool) ([]byte, error){
	keyIV := pbkdf2.Key(password, salt, 1000, AES256KeySize + aes.BlockSize, sha1.New)
	key, iv := keyIV[:AES256KeySize], keyIV[AES256KeySize:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	if encrypt {
		mode = cipher.NewCBCEncrypter(block, iv)
	}
	cipheredData := make([]byte, len(data))
	mode.CryptBlocks(cipheredData, data)
	return cipheredData, nil
}

func CipherSize(size int) int {
	const blocksize = 16
	return size + blocksize - size % blocksize
}
