// Package auth provides functionallity for authentication and authorisation
package auth

import (
	"mpt_data/database"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"

	"mpt_data/helper/errors"
)

func CreateUser(user apiModel.UserLogin) error {
	if user.Password == "" || user.Username == "" {
		return errors.ErrUserNotComplete
	}

	hash, err := hash(user.Password)
	if err != nil {
		return err
	}

	dbUser := dbModel.User{Username: user.Username, Hash: hash}
	if err := addUser(dbUser); err != nil {
		return err
	}
	return nil
}

func ChangePassword(userID uint, password string) error {
	if password == "" {
		return errors.ErrUserNotComplete
	}

	hash, err := hash(password)
	if err != nil {
		return err
	}

	if err :=
		database.DB.
			Table("users").
			Where("id = ?", userID).
			Update("hash", hash).Error; err != nil {
		return err
	}

	return nil
}

// AddUser creates a new user in the database
func addUser(user dbModel.User) error {
	db := database.DB.Begin()
	defer db.Rollback()

	if user.Username == "" || user.Hash == "" {
		return errors.ErrUserdataNotComplete
	}

	if rows := db.Find(&dbModel.User{}, "username = ?", user.Username).RowsAffected; rows != 0 {
		return errors.ErrUserAlreadyExists
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	if err := db.Commit().Error; err != nil {
		return err
	}
	return nil
}
