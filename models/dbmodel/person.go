package dbmodel

import (
	"encoding/json"

	"gorm.io/gorm"
)

// TODO: possibility to decide between people with same name
type Person struct {
	gorm.Model `json:"-"`
	ID         uint
	GivenName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
}

type PersonAbsence struct {
	gorm.Model `json:"-"`
	ID         uint
	MeetingID  uint    `gorm:"not null;index:personAbsence,unique" json:"-"`
	PersonID   uint    `gorm:"not null;index:personAbsence,unique" json:"-"`
	Meeting    Meeting `gorm:"ForeignKey:MeetingID"`
	Person     Person  `gorm:"ForeignKey:PersonID"`
}

type PersonRecurringAbsence struct {
	gorm.Model `json:"-"`
	ID         uint
	Weekday    int  `gorm:"not null;index:personRecurringAbsence,unique"`
	PersonID   uint `gorm:"not null;index:personRecurringAbsence,unique" json:"-"`
}

type PersonTask struct {
	gorm.Model   `json:"-"`
	ID           uint
	PersonID     uint       `gorm:"not null;index:personTask,unique" json:"-"`
	TaskDetailID uint       `gorm:"not null;index:personTask,unique" json:"-"`
	TaskDetail   TaskDetail `gorm:"ForeignKey:PersonID"`
	Person       Person     `gorm:"ForeignKey:TaskDetailID"`
}

func (pt PersonTask) MarshalJSON() ([]byte, error) {
	type Alias PersonTask
	if pt.TaskDetail.ID == 0 && pt.Person.ID == 0 {
		return json.Marshal(&struct {
			Alias
			TaskDetail interface{} `json:",omitempty"`
			Person     interface{} `json:",omitempty"`
		}{
			Alias:      (Alias)(pt),
			TaskDetail: nil,
			Person:     nil,
		})
	} else if pt.TaskDetail.ID == 0 {
		return json.Marshal(&struct {
			Alias
			TaskDetail interface{} `json:",omitempty"`
		}{
			Alias:      (Alias)(pt),
			TaskDetail: nil,
		})
	} else if pt.Person.ID == 0 {
		return json.Marshal(&struct {
			Alias
			Person interface{} `json:",omitempty"`
		}{
			Alias:  (Alias)(pt),
			Person: nil,
		})
	}
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(pt),
	})
}
