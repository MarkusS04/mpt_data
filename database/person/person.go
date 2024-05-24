package person

import (
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
)

const packageName = "database.person"

// AddPerson adds a Person to database
func AddPerson(person *dbModel.Person) error {
	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Create(&person).Error; err != nil {
		db.Rollback()
		logging.LogError(packageName+".AddPerson", err.Error())
		return err
	}

	return nil
}

// UpdatePerson changes the given-/lastName
func UpdatePerson(person dbModel.Person) error {
	if person.ID == 0 {
		return errors.ErrIDNotSet
	}

	db := database.DB.Begin()
	defer db.Commit()

	if err :=
		db.Model(&dbModel.Person{}).
			Where("id = ?", person.ID).
			Update("given_name", person.GivenName).
			Update("last_name", person.LastName).Error; err != nil {
		db.Rollback()
		return err
	}

	return nil
}

// GetPerson return a list of all people
func GetPerson() (people []dbModel.Person, err error) {
	if err := database.DB.Find(&people).Error; err != nil {
		return nil, err
	}
	return people, err
}

// DeletePerson delete a person, id must be set
func DeletePerson(person dbModel.Person) error {
	if person.ID == 0 {
		return errors.ErrIDNotSet
	}
	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Unscoped().Delete(&person).Error; err != nil {
		return err
	}

	return nil
}
