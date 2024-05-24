package dbmodel

import (
	generalmodel "mpt_data/models/general"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model   `json:"-"`
	ID           uint
	PersonID     uint       `json:"-"`
	MeetingID    uint       `json:"-" gorm:"not null;index:planTaskMeeting,unique"`
	TaskDetailID uint       `json:"-" gorm:"not null;index:planTaskMeeting,unique"`
	Person       Person     `gorm:"ForeignKey:PersonID"`
	Meeting      Meeting    `gorm:"ForeignKey:MeetingID"`
	TaskDetail   TaskDetail `gorm:"ForeignKey:TaskDetailID"`
}

type PDF struct {
	ID       uint
	Name     string
	FilePath string
	generalmodel.Period
	DataChanged bool
}
