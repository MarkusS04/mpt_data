// Package person provides functions to execute CRUD on person
package person

import (
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"

	"gorm.io/gorm"
)

const packageName = "database.person"

// AddPerson adds a Person to database
func AddPerson(db *gorm.DB, person *dbModel.Person) error {
	if err := db.Create(&person).Error; err != nil {
		logging.LogError(packageName+".AddPerson", err.Error())
		return err
	}

	return nil
}

// UpdatePerson changes the given-/lastName
func UpdatePerson(db *gorm.DB, person *dbModel.Person) error {
	if person.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err :=
		db.Save(person).Error; err != nil {
		return err
	}

	return nil
}

// GetPerson return a list of all people
func GetPerson(db *gorm.DB) (people []dbModel.Person, err error) {
	if err := db.Find(&people).Error; err != nil {
		return nil, err
	}
	return people, err
}

// DeletePerson delete a person, id must be set
func DeletePerson(db *gorm.DB, person dbModel.Person) error {
	if person.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err :=
		db.Unscoped().
			Delete(&person).
			Error; err != nil {
		return err
	}

	return nil
}
