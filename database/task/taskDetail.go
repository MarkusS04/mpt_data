// Package task provides functions to manipulate tasks and their details in database
package task

import (
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AddTaskDetail adds taskDetails to a task
func AddTaskDetail(db *gorm.DB, task *dbModel.TaskDetail) error {
	if task.TaskID == 0 {
		return errors.ErrForeignIDNotSet
	}

	if err := db.Save(&task).Error; err != nil {
		zap.L().Error(generalmodel.DBSaveDataFailed, zap.Error(err))
		return err
	}

	return nil
}

// UpdateTaskDetail updates descr of taskDetail
func UpdateTaskDetail(db *gorm.DB, task dbModel.TaskDetail) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err := db.Save(&task).Error; err != nil {
		zap.L().Error(generalmodel.DBUpdateDataFailed, zap.Error(err))
		return err
	}

	return nil
}

// DeleteTaskDetail deletes a taskDetail
func DeleteTaskDetail(db *gorm.DB, task dbModel.TaskDetail) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err := db.Unscoped().Delete(&task).Error; err != nil {
		zap.L().Error(generalmodel.DBDeleteDataFailed, zap.Error(err))
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
			zap.L().Error(generalmodel.DBLoadDataFailed, zap.Error(err))
			return err
		}
	}

	if err := db.Model(&dbModel.PDF{}).Where("1=1").Update("data_changed", true).Error; err != nil {
		zap.L().Error(generalmodel.DBUpdateDataFailed, zap.Error(err))
		return err
	}

	return nil
}
