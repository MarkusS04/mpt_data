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

const packageName = "database.task"

// GetTask loads all tasks based on the given conditions
func GetTask(db *gorm.DB, conds ...interface{}) (tasks []dbModel.Task, err error) {
	err = db.Preload("TaskDetails",
		func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number NULLS LAST")
		}).
		Order("order_number NULLS LAST").
		Find(&tasks, conds...).Error
	return tasks, err
}

// AddTask adds a task with taskDetails
func AddTask(db *gorm.DB, task *dbModel.Task) error {
	if err := db.Create(task).Error; err != nil {
		zap.L().Error(generalmodel.DBSaveDataFailed, zap.Error(err))
		return err
	}

	return nil
}

// UpdateTask sets a new Descr for a task, nothing else is changed
func UpdateTask(db *gorm.DB, task dbModel.Task) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err :=
		db.Save(&task).
			Error; err != nil {
		zap.L().Error(generalmodel.DBUpdateDataFailed, zap.Error(err))
		return err
	}

	return nil
}

// OrderTask will update the orderNumber in DB for the given Tasks
func OrderTask(db *gorm.DB, tasks []apimodel.OrderTask) error {
	for _, t := range tasks {
		if err :=
			db.Model(&dbModel.Task{}).
				Where("id = ?", t.TaskID).
				Update("order_number", t.OrderNumber).
				Error; err != nil && err != gorm.ErrRecordNotFound {
			zap.L().Error(generalmodel.DBLoadDataFailed, zap.Error(err))
			return err
		}
	}

	if err := db.Model(&dbModel.PDF{}).Where("1=1").Update("data_changed", true).Error; err != nil {
		zap.L().Error(generalmodel.DBLoadDataFailed, zap.Error(err))
		return err
	}

	return nil
}

// DeleteTask deletes a task with all taskDetails
func DeleteTask(db *gorm.DB, task *dbModel.Task) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}

	if err := db.Unscoped().Select("TaskDetails").Delete(&task).Error; err != nil {
		db.Rollback()
		zap.L().Error(generalmodel.DBDeleteDataFailed, zap.Error(err))
		return err
	}

	return nil
}
