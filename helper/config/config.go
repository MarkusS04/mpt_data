// Package config provides configuration and secret management
package config

import (
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
		Use          bool
		Environment  string
		WorkspaceID  string
		ClientID     string
		ClientSecret string
		JWTKey       bool
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
		if Config.SECRETS.JWTKey {
			Config.API.JWTKey = getKeyFromSecretProvider("JWT_KEY")
		} else {
			Config.API.JWTKey = os.ExpandEnv(Config.API.JWTKey)
		}
	} else {
		Config.API.JWTKey = os.ExpandEnv(Config.API.JWTKey)
	}

	Config.Database.Path = os.ExpandEnv(Config.Database.Path)
	Config.Log.Path = os.ExpandEnv(Config.Log.Path)
	Config.PDF.Path = os.ExpandEnv(Config.PDF.Path)

	createDirIfNotExist(Config.Database.Path)
	createDirIfNotExist(Config.Log.Path)
	createDirIfNotExist(Config.PDF.Path)
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
