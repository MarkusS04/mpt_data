package dbmodel

import (
	"mpt_data/helper"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username []byte `gorm:"not null; uniqueIndex"`
	Hash     string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null" json:"-"`
}

// Encrypt encrypts user data
func (u *User) Encrypt() error {
	// hash, err := helper.EncryptData([]byte(u.Hash))
	// if err != nil {
	// 	return err
	// }
	// u.Hash = string(hash)

	username, err := helper.EncryptDataDeterministic([]byte(u.Username))
	if err != nil {
		return err
	}
	u.Username = username

	return nil
}

// EncryptedUsername returns the encryptedUsername
func (u *User) EncryptedUsername() ([]byte, error) {
	username, err := helper.EncryptDataDeterministic(u.Username)
	if err != nil {
		return nil, err
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
	// hash, err := helper.DecryptData([]byte(u.Hash))
	// if err != nil {
	// 	return err
	// }
	// u.Hash = string(hash)

	username, err := helper.DecryptDataDeterministic([]byte(u.Username))
	if err != nil {
		return err
	}
	u.Username = username

	return
}
