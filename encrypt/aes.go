package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// Advanced Encryption Standard, AES
// For 16,24,32 bit strings, corresponding to AES-128, AES-192, and AES-256 encryption methods, respectively
// PwdKey should be confidential content
var PwdKey = []byte("DLW--#HAPIJWOCNB")

// PKCS7 Fill
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	// Repeat() function is to copy and paste [] byte {byte (padding)} slices
	// and then merge them into a new byte slice to return
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// Encryption
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	// 1. Create encryption instance
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
