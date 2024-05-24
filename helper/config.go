package helper

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	DataBasePath string
	Log          struct {
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

	var conf config
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	conf.API.JWTKey = os.ExpandEnv(conf.API.JWTKey)

	conf.DataBasePath = os.ExpandEnv(conf.DataBasePath)
	conf.Log.Path = os.ExpandEnv(conf.Log.Path)
	conf.PDF.Path = os.ExpandEnv(conf.PDF.Path)

	createDirIfNotExist(conf.DataBasePath)
	createDirIfNotExist(conf.Log.Path)
	createDirIfNotExist(conf.PDF.Path)

	Config = conf
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
