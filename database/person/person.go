// Package person provides functions to execute CRUD on person
package person

import (
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const packageName = "database.person"

// AddPerson adds a Person to database
func AddPerson(db *gorm.DB, person *dbModel.Person) error {
	if err := db.Create(&person).Error; err != nil {
		zap.L().Error(generalmodel.DBUpdateDataFailed, zap.Error(err))
		return err
	}

	return nil
}

// UpdatePerson changes the given-/lastName
func UpdatePerson(db *gorm.DB, person *dbModel.Person) (err error) {
	if person.ID == 0 {
		return errors.ErrIDNotSet
	}

	err = db.Save(person).Error
	return
}

// GetPerson return a list of all people
func GetPerson(db *gorm.DB) (people []dbModel.Person, err error) {
	if err := db.Find(&people).Error; err != nil {
		return nil, err
	}
	return
}

// DeletePerson delete a person, id must be set
func DeletePerson(db *gorm.DB, person dbModel.Person) (err error) {
	if person.ID == 0 {
		return errors.ErrIDNotSet
	}

	err = db.Unscoped().Delete(&person).Error

	return
}
