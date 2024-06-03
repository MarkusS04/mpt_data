// Package dbmodel provides all structs for databse ORM
package dbmodel

import (
	"encoding/json"
	"mpt_data/helper"
	"mpt_data/helper/errors"

	"gorm.io/gorm"
)

// Task is the primary group for tasks to be sorted in
type Task struct {
	gorm.Model  `json:"-"`
	ID          uint
	Descr       string       `gorm:"not null;uniqueIndex"`
	TaskDetails []TaskDetail `gorm:"ForeignKey:TaskID" json:",omitempty"`
	OrderNumber uint
}

// BeforeCreate hook for gorm
func (t *Task) BeforeCreate(db *gorm.DB) (err error) {
	if t.Descr == "" {
		return errors.ErrTaskDescrNotSet
	}

	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Descr == t.Descr {
			return errors.ErrTaskAlreadyExists
		}
	}

	descr, err := helper.EncryptDataToBase64(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = descr
	return
}

// AfterCreate hook for gorm
func (t *Task) AfterCreate(_ *gorm.DB) (err error) {
	descr, err := helper.DecryptDataFromBase64(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = string(descr)
	return
}

// BeforeUpdate hook for gorm
func (t *Task) BeforeUpdate(db *gorm.DB) (err error) {
	return t.BeforeCreate(db)
}

// AfterUpdate hook for gorm
func (t *Task) AfterUpdate(db *gorm.DB) (err error) {
	return t.AfterCreate(db)
}

// AfterFind hook for gorm
func (t *Task) AfterFind(_ *gorm.DB) (err error) {
	descr, err := helper.DecryptDataFromBase64(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = string(descr)
	return
}

// TaskDetail is a description for a specific Task, so it has always an associated Task
type TaskDetail struct {
	gorm.Model  `json:"-"`
	ID          uint
	Descr       string `gorm:"not null;index:taskDetailsUnique,unique"`
	TaskID      uint   `gorm:"not null;index:taskDetailsUnique,unique"`
	Task        Task   `json:",omitempty" gorm:"foreignkey:TaskID"`
	OrderNumber uint
}

// BeforeCreate hook for gorm
func (t *TaskDetail) BeforeCreate(db *gorm.DB) (err error) {
	if t.Descr == "" {
		return errors.ErrTaskDescrNotSet
	}

	var tasks []TaskDetail
	if err := db.Where("task_id = ?", t.TaskID).Find(&tasks).Error; err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Descr == t.Descr {
			return errors.ErrTaskDetailAlreadyExists
		}
	}

	descr, err := helper.EncryptDataToBase64(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = descr
	return
}

// AfterCreate hook for gorm
func (t *TaskDetail) AfterCreate(_ *gorm.DB) (err error) {
	descr, err := helper.DecryptDataFromBase64(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = string(descr)
	return
}

// BeforeUpdate hook for gorm
func (t *TaskDetail) BeforeUpdate(db *gorm.DB) (err error) {
	return t.BeforeCreate(db)
}

// AfterUpdate hook for gorm
func (t *TaskDetail) AfterUpdate(db *gorm.DB) (err error) {
	return t.AfterCreate(db)
}

// AfterFind hook for gorm
func (t *TaskDetail) AfterFind(_ *gorm.DB) (err error) {
	descr, err := helper.DecryptDataFromBase64(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = string(descr)
	return
}

// MarshalJSON marshals an TaskDetail into json
func (t TaskDetail) MarshalJSON() ([]byte, error) {
	type Alias TaskDetail
	if t.Task.ID == 0 {
		return json.Marshal(&struct {
			Alias
			Task interface{} `json:",omitempty"`
		}{
			Alias: (Alias)(t),
			Task:  nil,
		})
	}
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(t),
	})
}
