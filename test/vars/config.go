package vars

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/helper/config"
	"os"
)

func PrepareConfig() {
	if err := database.Connect(GetDbPAth()); err != nil {
		panic(err)
	}
	if err := os.Chdir("../.."); err != nil {
		fmt.Println(err)
	}
	// Load the config
	config.LoadConfig()
}
