// Package helper provides basic functionality
package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
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
func EncryptData(data []byte) ([]byte, error) {
	block, err := createCipherBlock()
	if err != nil {
		return nil, err
	}

	iv, err := generateRandomIV(block.BlockSize())
	if err != nil {
		return nil, err
	}

	return encrypt(block, iv, data), nil
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

// EncryptDataDeterministic encrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management, but uses a hash as IV
func EncryptDataDeterministic(data []byte) ([]byte, error) {
	block, err := createCipherBlock()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	iv := hash[:block.BlockSize()]

	return encrypt(block, iv, data), nil
}

// DecryptDataDeterministic decrypts data based on key stores in Config.Database.EncryptionKey, loaded from config or secret management, but uses a hash as IV
func DecryptDataDeterministic(data []byte) ([]byte, error) {
	return DecryptData(data)
}
