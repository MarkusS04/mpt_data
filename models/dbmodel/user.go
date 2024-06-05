// Package dbmodel provides all structs for databse ORM
package dbmodel

import (
	"mpt_data/helper"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null; uniqueIndex"`
	Hash     string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null" json:"-"`
}

// Encrypt encrypts user data
func (u *User) Encrypt() error {
	username, err := helper.EncryptDataDeterministicToBase64([]byte(u.Username))
	if err != nil {
		return err
	}
	u.Username = username

	return nil
}

// EncryptedUsername returns the encryptedUsername
func (u *User) EncryptedUsername() (string, error) {
	username, err := helper.EncryptDataDeterministicToBase64(u.Username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// BeforeCreate encryptes data in Database
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	return u.Encrypt()
}

// BeforeUpdate encryptes data in Database
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	return u.Encrypt()
}

// AfterFind decryptes data from Database
func (u *User) AfterFind(_ *gorm.DB) (err error) {
	return u.Decrypt()
}

// Decrypt decryptes user data
func (u *User) Decrypt() (err error) {
	username, err := helper.DecryptDataDeterministicFromBase64(u.Username)
	if err != nil {
		return err
	}
	u.Username = string(username)

	return
}
