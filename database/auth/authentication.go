// Package auth provides functionallity for authentication and authorisation
package auth

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Login tests user-credentials and if succesfull returns an JWT-Token
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
		Username: user.Username,
	}
	if err := userDb.Encrypt(); err != nil {
		return nil, err
	}

	if err := db.Where("username = ?", userDb.Username).First(&userDb).Error; err != nil {
		userDb.Decrypt()
		zap.L().Warn(generalmodel.UserInvalidLogin, zap.Error(err), zap.String("username", userDb.Username))
		return nil, errors.ErrInvalidAuth
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Hash), []byte(user.Password)); err != nil {
		zap.L().Warn(generalmodel.UserInvalidLogin, zap.Error(err))
		return nil, errors.ErrInvalidAuth
	}

	return &userDb, nil
}
