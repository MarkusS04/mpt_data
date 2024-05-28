package auth

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"

	"golang.org/x/crypto/bcrypt"
)

func Login(user apiModel.UserLogin) (string, error) {
	userDb, err := validateUser(user)
	if err != nil {
		return "", err
	}

	fmt.Println(userDb.ID)

	token, err := generateJWT(*userDb)
	if err != nil {
		return "", err
	}

	return token, nil
}

func hash(password string) ([]byte, error) {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return byteHash, nil
}

func validateUser(user apiModel.UserLogin) (*dbModel.User, error) {
	db := database.DB.Begin()
	defer db.Rollback()

	userDb := dbModel.User{
		Username: []byte(user.Username),
	}
	if err := userDb.Encrypt(); err != nil {
		return nil, err
	}

	if err := db.Where("username = ?", userDb.Username).First(&userDb).Error; err != nil {
		return nil, errors.ErrInvalidAuth
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Hash), []byte(user.Password)); err != nil {
		logging.LogWarning("database.auth.validateUser", err.Error())
		return nil, errors.ErrInvalidAuth
	}

	return &userDb, nil
}
