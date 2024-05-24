package task

import (
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
)

func AddTaskDetail(task *dbModel.TaskDetail) error {
	if task.TaskID == 0 {
		return errors.ErrForeignIDNotSet
	}
	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Save(&task).Error; err != nil {
		db.Rollback()
		logging.LogError(packageName+"AddTaskDetail", err.Error())
		return err
	}

	return nil
}

func UpdateTaskDetail(task dbModel.TaskDetail) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}

	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Model(&dbModel.TaskDetail{}).Where("id = ?", task.ID).Update("descr", task.Descr).Error; err != nil {
		db.Rollback()
		logging.LogError(packageName+".UpdateTaskDetail", err.Error())
		return err
	}

	return nil
}

func DeleteTaskDetail(task dbModel.TaskDetail) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}
	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Unscoped().Delete(&task).Error; err != nil {
		db.Rollback()
		logging.LogError(packageName+".DeleteTaskDetail", err.Error())
		return err
	}
	return nil
}
