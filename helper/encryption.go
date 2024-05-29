// Package helper provides basic functionality
package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"mpt_data/helper/config"
)

func createCipherBlock() (cipher.Block, error) {
	return aes.NewCipher(config.Config.GetDBEncryptionKey())
}

func generateRandomIV(blockSize int) ([]byte, error) {
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	return iv, nil
}

func encrypt(block cipher.Block, iv, data []byte) []byte {
	ciphertext := make([]byte, len(iv)+len(data))
	copy(ciphertext, iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[len(iv):], data)

	return ciphertext
}

// EncryptData encrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management
func EncryptData[T ~[]byte | ~string](data T) ([]byte, error) {
	var byteData []byte

	switch v := any(data).(type) {
	case []byte:
		byteData = v
	case string:
		byteData = []byte(v)
	default:
		return nil, errors.New("unsupported data type")
	}

	block, err := createCipherBlock()
	if err != nil {
		return nil, err
	}

	iv, err := generateRandomIV(block.BlockSize())
	if err != nil {
		return nil, err
	}

	return encrypt(block, iv, padData(byteData, aes.BlockSize)), nil
}

// EncryptDataToBase64 calls EncryptData and then converts output to base64 string
func EncryptDataToBase64[T ~[]byte | ~string](data T) (string, error) {
	encdata, err := EncryptData(data)

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encdata), nil
}

// DecryptData decrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management
func DecryptData[T ~[]byte | ~string](data T) ([]byte, error) {
	var byteData []byte

	switch v := any(data).(type) {
	case []byte:
		byteData = v
	case string:
		byteData = []byte(v)
	default:
		return nil, errors.New("unsupported data type")
	}
	block, err := aes.NewCipher(config.Config.GetDBEncryptionKey())
	if err != nil {
		return nil, err
	}

	if len(byteData) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := byteData[:aes.BlockSize]
	byteData = byteData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(byteData, byteData)

	return unpadData(byteData), nil
}

// DecryptDataFromBase64 decodes base64 string, then calls DecryptData
func DecryptDataFromBase64(data string) ([]byte, error) {
	encdata, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return DecryptData(encdata)
}

// EncryptDataDeterministic encrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management, but uses a hash as IV
func EncryptDataDeterministic(data []byte) ([]byte, error) {
	block, err := createCipherBlock()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	iv := hash[:block.BlockSize()]

	return encrypt(block, iv, padData(data, aes.BlockSize)), nil
}

// DecryptDataDeterministic decrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management, but uses a hash as IV
func DecryptDataDeterministic(data []byte) ([]byte, error) {
	return DecryptData(data)
}

func padData(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func unpadData(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
