// Package person provides functions to execute CRUD on person
package person

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"

	"go.uber.org/zap"
)

// AddTaskToPerson adds tasks to a person
func AddTaskToPerson(personID uint, tasks []dbModel.TaskDetail) (personTask []dbModel.PersonTask, err error) {
	if personID == 0 {
		return nil, errors.ErrIDNotSet
	}
	for _, task := range tasks {
		if task.ID == 0 {
			return nil, errors.ErrIDNotSet
		}
		personTask = append(personTask, dbModel.PersonTask{PersonID: personID, TaskDetailID: task.ID})
	}
	db := database.DB.Begin()
	defer db.Commit()

	if err = db.Create(&personTask).Error; err != nil {
		db.Rollback()
		zap.L().Error(generalmodel.DBSaveDataFailed, zap.Error(err))
		return nil, err
	}
	return personTask, nil
}

// DeleteTaskFromPerson deletes tasks from a person
func DeleteTaskFromPerson(personID uint, tasks []dbModel.TaskDetail) error {
	if personID == 0 {
		return errors.ErrIDNotSet
	}
	db := database.DB.Begin()
	defer db.Commit()

	for _, task := range tasks {
		if task.ID == 0 {
			db.Rollback()
			return errors.ErrIDNotSet
		}
		if err :=
			db.Unscoped().
				Where("person_id = ?", personID).
				Where("task_detail_id = ?", task.ID).
				Delete(&dbModel.PersonTask{}).Error; err != nil {
			db.Rollback()
			zap.L().Error(generalmodel.DBDeleteDataFailed, zap.Error(err))
			return err
		}
	}
	return nil
}

// GetTaskOfPerson loads tasks assigned to a person
func GetTaskOfPerson(personID uint) (task []dbModel.Task, err error) {
	if personID == 0 {
		return nil, errors.ErrIDNotSet
	}

	db := database.DB
	if err := db.
		Preload("TaskDetails",
			db.Table("task_details").Where("id in (?)", db.Table("person_tasks").Where("person_id = ?", personID).Select("task_detail_id"))).
		Find(&task,
			db.Where("id in (?)", db.Table("task_details").Where("id in (?)", db.Table("person_tasks").Where("person_id = ?", personID).Select("task_detail_id")).Select("task_id"))).
		Error; err != nil {
		return nil, err
	}
	return task, nil
}

// GetPersonWithTask loads all people with tasks assigned
func GetPersonWithTask() (persons []dbModel.PersonTask, err error) {
	if err :=
		database.DB.Preload("TaskDetail.Task").
			Preload("TaskDetail").
			Preload("Person").
			Find(&persons).
			Error; err != nil {
		return nil, err
	}

	return persons, err
}
