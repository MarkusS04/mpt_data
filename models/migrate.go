package models

import (
	"encoding/base64"
	"fmt"
	"mpt_data/database"
	"mpt_data/database/auth"
	"mpt_data/database/logging"
	"mpt_data/helper"
	"mpt_data/helper/config"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	"mpt_data/models/dbmodel"
	"os"
	"regexp"

	"gorm.io/gorm"
)

func Init() {
	db := database.DB
	encryptExistingData()
	if err := db.AutoMigrate(
		&dbmodel.Log{},
		&dbmodel.User{},
		&dbmodel.Meeting{},
		&dbmodel.Task{},
		&dbmodel.TaskDetail{},
		&dbmodel.Person{},
		&dbmodel.PersonTask{},
		&dbmodel.PersonAbsence{},
		&dbmodel.PersonRecurringAbsence{},
		&dbmodel.Plan{},
		&dbmodel.PDF{},
	); err != nil {
		// kein Logging in DB verfügbar
		config.Config.Log.LevelDB = ^uint(0)
		logging.LogError("models.Init", err.Error())
		os.Exit(1)
	}
	if err := auth.CreateUser(apimodel.UserLogin{Username: "admin", Password: "admin"}); err != nil && err != errors.ErrUserAlreadyExists {
		fmt.Println(err)
	}
	fmt.Println("database initialized")
}

// will be dropped in future versions, it is that version 2 to 3 does not break the entire database system
func encryptExistingData() {

	if testDataAlreadyEncrypted() {
		return
	}
	database.DB.Exec("drop table logs")

	database.DB.Transaction(
		func(tx *gorm.DB) error {
			var tags []dbmodel.Tag
			tx.Session((&gorm.Session{SkipHooks: true})).Find(&tags)
			if len(tags) == 0 {
				return nil
			}
			if err := tx.Save(&tags).Error; err != nil {
				return err
			}
			return nil
		},
	)

	database.DB.Transaction(
		func(tx *gorm.DB) error {
			var people []dbmodel.Person
			tx.Session((&gorm.Session{SkipHooks: true})).Find(&people)
			if len(people) == 0 {
				return nil
			}
			if err := tx.Save(&people).Error; err != nil {
				return err
			}
			return nil
		},
	)

	database.DB.Transaction(
		func(tx *gorm.DB) error {
			var tasks []dbmodel.TaskDetail
			tx.Session((&gorm.Session{SkipHooks: true})).Find(&tasks)
			if len(tasks) == 0 {
				return nil
			}
			for i := range tasks {
				var err error
				tasks[i].Descr, err = helper.EncryptData(tasks[i].Descr)
				if err != nil {
					panic(err)
				}
			}
			if err := tx.Session((&gorm.Session{SkipHooks: true})).Save(&tasks).Error; err != nil {
				return err
			}
			return nil
		},
	)

	database.DB.Transaction(
		func(tx *gorm.DB) error {
			var tasks []dbmodel.Task
			tx.Session((&gorm.Session{SkipHooks: true})).Find(&tasks)
			if len(tasks) == 0 {
				return nil
			}
			for i := range tasks {
				var err error
				tasks[i].Descr, err = helper.EncryptData(tasks[i].Descr)
				if err != nil {
					panic(err)
				}
			}
			if err := tx.Session((&gorm.Session{SkipHooks: true})).Save(&tasks).Error; err != nil {
				return err
			}
			return nil
		},
	)

	database.DB.Transaction(
		func(tx *gorm.DB) error {
			var users []dbmodel.User
			tx.Session((&gorm.Session{SkipHooks: true})).Find(&users)
			if len(users) == 0 {
				return nil
			}
			if err := tx.Save(&users).Error; err != nil {
				return err
			}
			return nil
		},
	)
}

// test some tables if data is already encrypted
func testDataAlreadyEncrypted() bool {
	var people []dbmodel.Person
	if rows := database.DB.Session((&gorm.Session{SkipHooks: true})).Find(&people).RowsAffected; rows != 0 {
		for _, person := range people {
			if !isBase64(person.GivenName) {
				return false
			}
		}
	}

	var logs []dbmodel.Log
	if rows := database.DB.Session((&gorm.Session{SkipHooks: true})).Find(&logs).RowsAffected; rows != 0 {
		for _, log := range logs {
			if !isBase64(log.Text) {
				return false
			}
		}
	}

	var tasks []dbmodel.Task
	if rows := database.DB.Session((&gorm.Session{SkipHooks: true})).Find(&tasks).RowsAffected; rows != 0 {
		for _, task := range tasks {
			if !isBase64(task.Descr) {
				return false
			}
		}
	}
	return true
}

func isBase64(s string) bool {
	// Check if the length is a multiple of 4
	if len(s)%4 != 0 {
		return false
	}

	// Check if the string contains only valid Base64 characters
	base64Pattern := regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	if !base64Pattern.MatchString(s) {
		return false
	}

	// Try to decode the string
	_, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}

	return true
}
