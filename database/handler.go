package database

import (
	"mpt_data/helper/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func Connect(dbPath string) error {
	var conf gorm.Config
	if !config.Config.Log.GormOutputEnabled {
		conf.Logger = logger.Default.LogMode(logger.Silent)
	}
	db, err := gorm.Open(
		sqlite.Open(dbPath+"/datbase.db"),
		&conf,
	)
	if err != nil {
		return err
	}

	// Setze Verbindungspool-Einstellungen
	sqldb, _ := db.DB()
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)

	DB = db
	return nil
}
