package task

import (
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"

	"gorm.io/gorm"
)

const packageName = "database.task"

// GetTask loads all tasks based on the given conditions
func GetTask(conds ...interface{}) (tasks []dbModel.Task, err error) {
	err = database.DB.
		Preload("TaskDetails", func(db *gorm.DB) *gorm.DB { return db.Order("order_number NULLS LAST") }).
		Order("order_number NULLS LAST").
		Find(&tasks, conds...).Error
	return tasks, err
}

// AddTask adds a task with taskDetails
func AddTask(task *dbModel.Task) error {
	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Create(task).Error; err != nil {
		db.Rollback()
		logging.LogError(packageName+".AddTask", err.Error())
		return err
	}

	return nil
}

// UpdateTask sets a new Descr for a task, nothing else is changed
func UpdateTask(task dbModel.Task) error {
	if task.ID == 0 {
		return errors.ErrIDNotSet
	}
	db := database.DB.Begin()
	defer db.Commit()

	if err := db.Model(&task).Where("id = ?", task.ID).Update("descr", task.Descr).Error; err != nil {
		db.Rollback()
		logging.LogError(packageName+".UpdateTask", err.Error())
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

// DeleteTask deletes a task with all taskDetails
func DeleteTask(task dbModel.Task) error {
	db := database.DB.Begin()
	defer db.Commit()

	if task.ID == 0 {
		if err := db.First(&task, "descr = ?", task.Descr).Error; err != nil {
			db.Rollback()
			logging.LogError(packageName+".DeleteTask", err.Error())
			return err
		}
	}

	if err := db.Unscoped().Select("TaskDetails").Delete(&task).Error; err != nil {
		db.Rollback()
		logging.LogError(packageName+".DeleteTask", err.Error())
		return err
	}

	return nil
}
