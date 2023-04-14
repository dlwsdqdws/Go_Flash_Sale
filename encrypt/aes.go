package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// Advanced Encryption Standard, AES
// Bidirectional encryption
// For 16,24,32 bit strings, corresponding to AES-128, AES-192, and AES-256 encryption methods, respectively

// PwdKey should be confidential content
var PwdKey = []byte("DLW--#HAPIJWOCNB")

// PKCS7 padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	// Repeat() function is to copy and paste [] byte {byte (padding)} slices
	// and then merge them into a new byte slice to return
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// Deleting padding strings
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("Invalid encrypted string！")
	} else {
		unpadding := int(origData[length-1])
		// Cut slices, remove padding bytes, and return clear text
		return origData[:(length - unpadding)], nil
	}
}

// Encryption
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	// 1. Create encryption algo instance
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 2. Get block(PwdKey) size
	blockSize := block.BlockSize()
	// 3. Fill in the data to ensure that the data length meets the requirements
	origData = PKCS7Padding(origData, blockSize)
	// 4. Using CBC encryption mode in AES encryption method
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	// 5. Perform encryption
	blocMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// AES decryption
func AesDeCrypt(encypted []byte, key []byte) ([]byte, error) {
	// 1. Create encryption algo instance
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 2. Get block(PwdKey) size
	blockSize := block.BlockSize()
	// 3. Create encryption client instance
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encypted))
	// 4. Decryption - the same function as encryption
	blockMode.CryptBlocks(origData, encypted)
	// 5. Remove padding string
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// Encrypt base64
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEcrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// Decrypt base64
func DePwdCode(pwd string) ([]byte, error) {
	// decryption base64
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	// AES decryption
	return AesDeCrypt(pwdByte, PwdKey)
}
