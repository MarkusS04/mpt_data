package main

import (
	"fmt"
	"mpt_data/api"
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/helper"
	"mpt_data/models"
	"net/http"
	"os"
)

// @title						MPT
// @version					1
// @description				Meeting Planning Tool API
// @license.name				MIT
// @license.url				https://opensource.org/license/MIT
// @BasePath					/api/v1
// @SecurityDefinitions.apiKey	ApiKeyAuth
// @In							header
// @Name						Authorization
func main() {
	helper.LoadConfig()
	if err := database.Connect(helper.Config.DataBasePath); err != nil {
		fmt.Println("could not connect to database:", err)
		os.Exit(1)
	}
	models.Init()

	if err := http.ListenAndServe(":"+helper.Config.API.Port, api.PrepareServer()); err != nil {
		logging.LogError("main", "Failed to start API")
		os.Exit(1)
	}
}
