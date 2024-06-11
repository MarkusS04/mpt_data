package models

import (
	"mpt_data/database"
	"mpt_data/database/auth"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	"mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"os"

	"go.uber.org/zap"
)

func Init() {
	db := database.DB

	if err := db.AutoMigrate(
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
		zap.L().Error(generalmodel.DBMigrationFailed, zap.Error(err))
		os.Exit(1)
	}

	if err := auth.CreateUser(apimodel.UserLogin{Username: "admin", Password: "admin"}); err != nil && err != errors.ErrUserAlreadyExists {
		zap.L().Error(generalmodel.UserCreationFailed, zap.Error(err))
	}

	zap.L().Info(generalmodel.DBMigrated)
}
