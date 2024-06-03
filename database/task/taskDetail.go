package task

import (
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"

	"gorm.io/gorm"
)

func AddTaskDetail(db *gorm.DB, task *dbModel.TaskDetail) error {
	if task.TaskID == 0 {
		return errors.ErrForeignIDNotSet
	}

	if err := db.Save(&task).Error; err != nil {
		logging.LogError(packageName+"AddTaskDetail", err.Error())
		return err
	}

	return nil
}

func UpdateTaskDetail(db *gorm.DB, task dbModel.TaskDetail) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err := db.Save(&task).Error; err != nil {
		logging.LogError(packageName+".UpdateTaskDetail", err.Error())
		return err
	}

	return nil
}

func DeleteTaskDetail(db *gorm.DB, task dbModel.TaskDetail) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err := db.Unscoped().Delete(&task).Error; err != nil {
		logging.LogError(packageName+".DeleteTaskDetail", err.Error())
		return err
	}
	return nil
}

// OrderTaskDetail will update the orderNumber in DB for the given Tasks
func OrderTaskDetail(db *gorm.DB, tasks []apimodel.OrderTaskDetail, taskID uint) error {
	for _, t := range tasks {
		if err :=
			db.Model(&dbModel.TaskDetail{}).
				Where("id = ?", t.TaskDetailID).
				Where("task_id = ?", taskID).
				Update("order_number", t.OrderNumber).
				Error; err != nil && err != gorm.ErrRecordNotFound {
			logging.LogError(packageName+".OrderTask", err.Error())
			return err
		}
	}

	if err := db.Model(&dbModel.PDF{}).Where("1=1").Update("data_changed", true).Error; err != nil {
		logging.LogError(packageName+".OrderTask", err.Error())
		return err
	}

	return nil
}
