// Package absence provides CRUD for absences
package absence

import (
	"mpt_data/helper/errors"
	"mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AddRecurringAbsence stores recurring absence for person
func AddRecurringAbsence(absences []*dbmodel.PersonRecurringAbsence, db *gorm.DB) error {
	if err := db.Save(&absences).Error; err != nil {
		zap.L().Error(generalmodel.DBSaveDataFailed, zap.Error(err))
		return err
	}
	return nil
}

// DeleteRecurringAbsence deletes recurring absence for person
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
			zap.L().Error(generalmodel.DBDeleteDataFailed, zap.Error(err))
			return err
		}
	}
	return nil
}

// GetRecurringAbsence Loads deletes recurring absence for person
func GetRecurringAbsence(personID uint, db *gorm.DB) (absences []dbmodel.PersonRecurringAbsence, err error) {
	if err :=
		db.
			Where("person_id = ?", personID).
			Order("weekday asc").
			Find(&absences).Error; err != nil {
		zap.L().Error(generalmodel.DBLoadDataFailed, zap.Error(err), zap.Uint("person_id", personID))
		return nil, err
	}

	return absences, nil
}
