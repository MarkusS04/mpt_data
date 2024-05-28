package main

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/database/auth"
	"mpt_data/helper/config"
	"mpt_data/models"
	"mpt_data/test/vars"
	"os"
)

// main must run to prepare the test environment
func main() {
	config.LoadConfig()
	dbPath := vars.GetDbPAth()
	if err := os.RemoveAll(dbPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := database.Connect(vars.GetDbPAth()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	models.Init()
	prepareTestData()
}

func prepareTestData() {
	os.Setenv("MPT_JWT_KEY", "a_super_secret_key_for_testing")
	if err := auth.CreateUser(vars.UserAPI); err != nil {
		fmt.Println("Error preparing testdata")
		os.Exit(1)
	}
}
