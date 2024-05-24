package database

import (
	"mpt_data/helper"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func Connect(dbPath string) error {
	var config gorm.Config
	if !helper.Config.Log.GormOutputEnabled {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	db, err := gorm.Open(
		sqlite.Open(dbPath+"/datbase.db"),
		&config,
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
