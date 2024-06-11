package main

import (
	"fmt"
	"mpt_data/api"
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/database/plan"
	"mpt_data/helper/config"
	"mpt_data/models"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
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
	// Prepare Database
	config.LoadConfig()
	if err := database.Connect(config.Config.Database.Path); err != nil {
		fmt.Println("could not connect to database:", err)
		os.Exit(1)
	}
	models.Init()

	// Delete PDFs after x months every month
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() {
		logging.LogInfo("main.cron", "Start execution of pdf autoremoval")
		plan.PDFAutoRemoval(database.DB, 3)
		logging.LogInfo("main.cron", "End execution of pdf autoremoval")
	})
	go c.Start()

	// Start WebServer
	if err := http.ListenAndServe(":"+config.Config.API.Port, api.PrepareServer()); err != nil {
		logging.LogError("main", "Failed to start API")
		os.Exit(1)
	}
}
