package main

import (
	"fmt"
	"mpt_data/api"
	"mpt_data/database"
	"mpt_data/database/plan"
	"mpt_data/helper/config"
	"mpt_data/models"
	generalmodel "mpt_data/models/general"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func init() {
	// Prepare Database
	config.LoadConfig()
	if err := database.Connect(config.Config.Database.Path); err != nil {
		fmt.Println("could not connect to database:", err)
		os.Exit(1)
	}
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{config.Config.Log.Path + "/log.log"}
	zap.ReplaceGlobals(zap.Must(cfg.Build()))
	models.Init()

	// Delete PDFs after x months every month
	c := cron.New()
	c.AddFunc("0 0 1 * *", func() {
		zap.L().Info(generalmodel.StartExecPDFAutoremoval)
		plan.PDFAutoRemoval(database.DB, 3)
		zap.L().Info(generalmodel.EndExecPDFAutoremoval)
	})
	go c.Start()
}

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
	// Start WebServer
	if err := http.ListenAndServe(":"+config.Config.API.Port, api.PrepareServer()); err != nil {
		zap.L().Error(generalmodel.APIStartFailed, zap.Error(err))
		os.Exit(1)
	}
}
