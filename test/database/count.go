package database

import (
	"fmt"
	"mpt_data/database"

	"gorm.io/gorm"
)

func CountEntries(table interface{}, conds ...interface{}) (count int64) {
	db := database.DB.Begin()
	defer db.Rollback()

	result := db.Find(&table, conds...)

	if result.Error != nil {
		fmt.Println("error: ", result.Error)
	}

	result.Count(&count)
	return count
}

func CountEntriesDB(db *gorm.DB, table interface{}, conds ...interface{}) (count int64) {

	result := db.Find(&table, conds...)

	if result.Error != nil {
		fmt.Println("error: ", result.Error)
	}

	result.Count(&count)
	return count
}
