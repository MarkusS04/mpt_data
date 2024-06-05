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

// EncryptData calls EncryptData and then converts output to base64 string
func EncryptData[T ~[]byte | ~string](data T) (string, error) {
	var byteData []byte

	switch v := any(data).(type) {
	case []byte:
		byteData = v
	case string:
		byteData = []byte(v)
	default:
		return "", errors.New("unsupported data type")
	}

	block, err := createCipherBlock()
	if err != nil {
		return "", err
	}

	iv, err := generateRandomIV(block.BlockSize())
	if err != nil {
		return "", err
	}

	encdata := encrypt(block, iv, padData(byteData, aes.BlockSize))
	return base64.StdEncoding.EncodeToString(encdata), nil
}

// DecryptData decodes base64 string, then calls DecryptData
func DecryptData(data string) ([]byte, error) {
	encdata, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(config.Config.GetDBEncryptionKey())
	if err != nil {
		return nil, err
	}

	if len(encdata) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := encdata[:aes.BlockSize]
	encdata = encdata[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encdata, encdata)

	return unpadData(encdata), nil
}

// EncryptDataDeterministicToBase64 encrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management, but uses a hash as IV
func EncryptDataDeterministicToBase64[T ~[]byte | ~string](data T) (string, error) {
	var byteData []byte
	switch v := any(data).(type) {
	case []byte:
		byteData = v
	case string:
		byteData = []byte(v)
	default:
		return "", errors.New("unsupported data type")
	}
	block, err := createCipherBlock()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(byteData)
	iv := hash[:block.BlockSize()]

	enc := encrypt(block, iv, padData(byteData, aes.BlockSize))
	return base64.StdEncoding.EncodeToString(enc), err
}

// DecryptDataDeterministicFromBase64 decrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management, but uses a hash as IV
func DecryptDataDeterministicFromBase64(data string) ([]byte, error) {
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
