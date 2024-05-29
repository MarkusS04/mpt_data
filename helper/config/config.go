// Package config provides configuration and secret management
package config

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config stores config based on a config file, and on requests made to secret management platform
type config struct {
	Database struct {
		Path          string
		EncryptionKey string
	}
	Log struct {
		Path              string
		LevelDB           uint
		GormOutputEnabled bool
	}
	API struct {
		Port                   string
		JWTKey                 string
		UseSwagger             bool
		AuthenticationRequired bool
	}
	PDF struct {
		Path string
	}

	SECRETS struct {
		Use             bool
		Environment     string
		WorkspaceID     string
		ClientID        string
		ClientSecret    string
		JWTKey          bool
		DBEncryptionKey bool
	}
}

func (conf *config) GetDBEncryptionKey() []byte {
	key, err := base64.StdEncoding.DecodeString(conf.Database.EncryptionKey)
	if err != nil {
		panic(fmt.Sprintf("please check your encryption key: %v", err))
	}
	return key
}
func (conf *config) testConfig() {
	if len(conf.API.JWTKey) == 0 {
		panic("JWT key is empty")
	}
	keyLen := len(conf.GetDBEncryptionKey())

	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		panic("invalid AES key length: must be 16, 24, or 32 bytes, and Base64-encoded")
	}
}

var Config config

// LoadConfig reads in the config file
// Creates all directorys if they not exits
func LoadConfig() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if err := viper.Unmarshal(&Config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	if Config.SECRETS.Use {
		loginToSecretProvider()
	}

	Config.API.JWTKey = getKey(Config.API.JWTKey, "JWT_KEY", Config.SECRETS.Use && Config.SECRETS.JWTKey)
	Config.Database.EncryptionKey = getKey(Config.Database.EncryptionKey, "DB_ENCRYPTION_KEY", Config.SECRETS.Use && Config.SECRETS.DBEncryptionKey)

	Config.Database.Path = os.ExpandEnv(Config.Database.Path)
	Config.Log.Path = os.ExpandEnv(Config.Log.Path)
	Config.PDF.Path = os.ExpandEnv(Config.PDF.Path)

	createDirIfNotExist(Config.Database.Path)
	createDirIfNotExist(Config.Log.Path)
	createDirIfNotExist(Config.PDF.Path)

	Config.testConfig()
}

func getKey(configValue string, secretName string, useSecrets bool) string {
	if useSecrets {
		return getKeyFromSecretProvider(secretName)
	}
	return os.ExpandEnv(configValue)
}

func createDirIfNotExist(path string) {
	// Check if the folder exists
	_, err := os.Stat(path)
	if err == nil {
		return
	}

	// If the folder does not exist, create it
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Println("could not create folders specified in config:", err)
			os.Exit(1)
		}
	}
}
