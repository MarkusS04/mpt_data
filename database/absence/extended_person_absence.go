package absence

import (
	"fmt"
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	"mpt_data/models/dbmodel"

	"gorm.io/gorm"
)

func AddRecurringAbsence(absences []*dbmodel.PersonRecurringAbsence, db *gorm.DB) error {
	if err := db.Save(&absences).Error; err != nil {
		logging.LogError(packageName+".AddRecurringAbsence", err.Error())
		return err
	}
	return nil
}

func DeleteRecurringAbsence(absences []dbmodel.PersonRecurringAbsence, db *gorm.DB) error {
	const funcName = packageName + ".DeleteRecurringAbsence"
	for _, absence := range absences {
		if absence.PersonID == 0 || absence.Weekday < 0 || absence.Weekday > 6 {
			return errors.ErrIDNotSet
		}
	}

	for _, absence := range absences {
		if err :=
			db.Where("weekday = ?", absence.Weekday).
				Where("person_id = ?", absence.PersonID).
				Unscoped().Delete(&absence).Error; err != nil {
			logging.LogError(funcName, err.Error())
			return err
		}
	}
	return nil
}

func GetRecurringAbsence(personID uint, db *gorm.DB) (absences []dbmodel.PersonRecurringAbsence, err error) {
	if err :=
		db.
			Where("person_id = ?", personID).
			Order("weekday asc").
			Find(&absences).Error; err != nil {
		logging.LogError(packageName+".GetRecurringAbsence", fmt.Sprintf("%s, person_id=%d", err.Error(), personID))
		return nil, err
	}

	return absences, nil
}
