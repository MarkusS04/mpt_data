// Package helper provides basic functionality
package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"mpt_data/helper/config"
)

// EncryptData encrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management
func EncryptData(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(config.Config.GetDBEncryptionKey())
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// DecryptData decrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management
func DecryptData(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(config.Config.GetDBEncryptionKey())
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
