package meeting

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
)

func GetMeetings(period generalmodel.Period) (meetings []dbModel.Meeting, err error) {
	db := database.DB.Begin()
	defer db.Rollback()

	err = db.Order("date asc").Where("date between ? and ?", period.StartDate, period.EndDate).
		Find(&meetings).Error

	return meetings, err
}

func AddMeetings(meetings []dbModel.Meeting) (err error) {

	db := database.DB.Begin()
	defer db.Commit()

	result := db.Create(meetings)

	if result.Error != nil {
		db.Rollback()
		return result.Error
	}

	if result.RowsAffected != int64(len(meetings)) {
		err = errors.ErrNotAllMeetingsCreated
	}

	return err
}

func UpdateMeeting(meeting dbModel.Meeting) (err error) {
	db := database.DB.Begin()
	defer db.Commit()

	result := db.Save(&meeting)
	if result.Error != nil {
		err = result.Error
		db.Rollback()
	}

	return err
}

func DeleteMeeting(meeting dbModel.Meeting) (err error) {
	db := database.DB.Begin()
	defer db.Commit()

	result := db.Unscoped().Delete(meeting)
	if result.Error != nil {
		db.Rollback()
		return result.Error
	}

	if result.RowsAffected != 1 {
		db.Rollback()
		return errors.ErrMeetingNotDeleted
	}

	return nil
}
