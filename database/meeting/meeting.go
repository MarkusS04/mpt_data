// Package meeting provides functions to manipulate meetings
package meeting

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"

	"gorm.io/gorm"
)

// GetMeetings returns a list of all meetings in the passed period
func GetMeetings(period generalmodel.Period) (meetings []dbModel.Meeting, err error) {
	db := database.DB.Begin()
	defer db.Rollback()

	err = db.Preload("Tag").Order("date asc").Where("date between ? and ?", period.StartDate, period.EndDate).
		Find(&meetings).Error

	return meetings, err
}

// AddMeetings creates all passed meetings in db, if doesnt exists already
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

// CreateTag creates a tag and links it to the meeting
func CreateTag(db *gorm.DB, meetingID uint, tag dbModel.Tag) error {
	if err := db.Create(&tag).Error; err != nil {
		return err
	}

	result := db.Model(&dbModel.Meeting{}).
		Where("id = ?", meetingID).
		Where("tag_id is null or tag_id = 0").
		Update("tag_id", tag.ID)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return errors.ErrMeetingTagAlreadySet
	}

	if err :=
		db.Model(&dbModel.PDF{}).
			Where("start_date <= (select date from meetings where id = ?)", meetingID).
			Where("end_date >= (select date from meetings where id = ?)", meetingID).
			Update("data_changed", true).
			Error; err != nil {
		return err
	}

	return nil
}

// DeleteTag deletes a tag from meeting
func DeleteTag(db *gorm.DB, meetingID uint) error {

	if err :=
		db.Where("id = (select tag_id from meetings where id = ?)", meetingID).
			Delete(&dbModel.Tag{}).
			Error; err != nil {
		return err
	}

	if err :=
		db.Model(&dbModel.Meeting{}).
			Where("id = ?", meetingID).
			Update("tag_id", nil).Error; err != nil {
		return err
	}

	if err :=
		db.Model(&dbModel.PDF{}).
			Where("start_date <= (select date from meetings where id = ?)", meetingID).
			Where("end_date >= (select date from meetings where id = ?)", meetingID).
			Update("data_changed", true).
			Error; err != nil {
		return err
	}

	return nil
}

// UpdateMeeting saves the passed meeting or creates if not exists
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

// DeleteMeeting deletes the passed meeting from db
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
