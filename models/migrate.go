package models

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/database/auth"
	"mpt_data/database/logging"
	"mpt_data/helper/config"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"os"
)

func Init() {
	db := database.DB
	if err := db.AutoMigrate(
		&dbModel.Log{},
		&dbModel.User{},
		&dbModel.Meeting{},
		&dbModel.Task{},
		&dbModel.TaskDetail{},
		&dbModel.Person{},
		&dbModel.PersonTask{},
		&dbModel.PersonAbsence{},
		&dbModel.PersonRecurringAbsence{},
		&dbModel.Plan{},
		&dbModel.PDF{},
	); err != nil {
		// kein Logging in DB verfügbar
		config.Config.Log.LevelDB = ^uint(0)
		logging.LogError("models.Init", err.Error())
		os.Exit(1)
	}
	if err := auth.CreateUser(apiModel.UserLogin{Username: "admin", Password: "admin"}); err != nil && err != errors.ErrUserAlreadyExists {
		fmt.Println(err)
	}
	fmt.Println("database initialized")
}
