package auth

import (
	"mpt_data/database"
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

	token, err := generateJWT(*userDb)
	if err != nil {
		return "", err
	}

	return token, nil
}

func hash(password string) (string, error) {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(byteHash), nil
}

func validateUser(user apiModel.UserLogin) (*dbModel.User, error) {
	db := database.DB.Begin()
	defer db.Rollback()

	var userDb dbModel.User
	if err := db.Where("username = ?", user.Username).First(&userDb).Error; err != nil {
		return nil, errors.ErrInvalidAuth
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Hash), []byte(user.Password)); err != nil {
		return nil, errors.ErrInvalidAuth
	}

	return &userDb, nil
}
