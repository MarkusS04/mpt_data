// Package absence provides CRUD for absences
package absence

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"sort"

	"go.uber.org/zap"
)

const packageName = "database.absence"

// AddAbsence stores absence
func AddAbsence(absences []dbModel.PersonAbsence) error {
	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Save(&absences).Error; err != nil {
		db.Rollback()
		zap.L().Error(generalmodel.DBSaveDataFailed, zap.Error(err))
		return err
	}
	return nil
}

// DeleteAbsence deletes absence
func DeleteAbsence(absences []dbModel.PersonAbsence) error {
	for _, absence := range absences {
		if absence.MeetingID == 0 || absence.PersonID == 0 {
			return errors.ErrIDNotSet
		}
	}
	db := database.DB.Begin()
	defer db.Commit()

	for _, absence := range absences {
		if err :=
			db.Where("meeting_id = ?", absence.MeetingID).
				Where("person_id = ?", absence.PersonID).
				Unscoped().Delete(&absence).Error; err != nil {
			db.Rollback()
			zap.L().Error(generalmodel.DBDeleteDataFailed, zap.Error(err))
			return err
		}
	}
	return nil
}

// GetAbsenceMeeting loads absence for specific meeting
func GetAbsenceMeeting(meetingID uint) (people []dbModel.Person, err error) {
	var meeting []dbModel.PersonAbsence
	if err :=
		database.DB.
			Preload("Person").
			Preload("Meeting").
			Where("meeting_id = ?", meetingID).
			Find(&meeting).Error; err != nil {
		zap.L().Error(generalmodel.DBLoadDataFailed, zap.Error(err), zap.Uint("meeting_id", meetingID))
		return nil, err
	}

	for _, entry := range meeting {
		people = append(people, entry.Person)
	}

	return people, nil
}

// GetAbsencePerson loads absence for specific person
func GetAbsencePerson(personID uint, period generalmodel.Period) (meetings []dbModel.Meeting, err error) {
	var absence []dbModel.PersonAbsence
	if err :=
		database.DB.
			Preload("Person").
			Preload("Meeting").
			Where("meeting_id IN (?)",
				database.DB.
					Table("meetings").
					Where("date BETWEEN ? AND ?", period.StartDate, period.EndDate).
					Select("id")).
			Where("person_id = ?", personID).
			Find(&absence).Error; err != nil {
		zap.L().Error(generalmodel.DBLoadDataFailed, zap.Error(err), zap.Uint("person_id", personID))

		return nil, err
	}

	sort.SliceStable(absence, func(i, j int) bool {
		return absence[i].Meeting.Date.Before(absence[j].Meeting.Date)
	})

	for _, entry := range absence {
		meetings = append(meetings, entry.Meeting)
	}

	return meetings, nil
}
