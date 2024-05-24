// Package dbmodel provides all structs for databse ORM
package dbmodel

import (
	"encoding/json"

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

// TaskDetail is a description for a specific Task, so it has always an associated Task
type TaskDetail struct {
	gorm.Model  `json:"-"`
	ID          uint
	Descr       string `gorm:"not null;index:taskDetailsUnique,unique"`
	TaskID      uint   `gorm:"not null;index:taskDetailsUnique,unique"`
	Task        Task   `json:",omitempty" gorm:"foreignkey:TaskID"`
	OrderNumber uint
}

// MarshalJSON marshals an TaskDetail into json
func (pt TaskDetail) MarshalJSON() ([]byte, error) {
	type Alias TaskDetail
	if pt.Task.ID == 0 {
		return json.Marshal(&struct {
			Alias
			Task interface{} `json:",omitempty"`
		}{
			Alias: (Alias)(pt),
			Task:  nil,
		})
	}
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(pt),
	})
}
