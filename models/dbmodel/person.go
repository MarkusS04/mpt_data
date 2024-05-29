// Package dbmodel provides all structs for databse ORM
package dbmodel

import (
	"encoding/json"
	"mpt_data/helper"
	"mpt_data/helper/errors"

	"gorm.io/gorm"
)

// TODO: possibility to decide between people with same name
// Person stores persons data
type Person struct {
	gorm.Model `json:"-"`
	ID         uint
	GivenName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
}

func (p *Person) encrypt() error {
	givenName, err := helper.EncryptDataToBase64(p.GivenName)
	if err != nil {
		return err
	}
	p.GivenName = givenName

	lastName, err := helper.EncryptDataToBase64(p.LastName)
	if err != nil {
		return err
	}
	p.LastName = lastName

	return nil
}

func (p *Person) decrypt() error {
	givenName, err := helper.DecryptDataFromBase64(p.GivenName)
	if err != nil {
		return err
	}
	p.GivenName = string(givenName)

	lastName, err := helper.DecryptDataFromBase64(p.LastName)
	if err != nil {
		return err
	}
	p.LastName = string(lastName)

	return nil
}

// BeforeCreate encryptes data in Database
func (p *Person) BeforeCreate(_ *gorm.DB) (err error) {
	if p.GivenName == "" || p.LastName == "" {
		return errors.ErrPersonMissingName
	}
	return p.encrypt()
}

// AfterCreate decryptes data after creation
func (p *Person) AfterCreate(_ *gorm.DB) (err error) {
	return p.decrypt()
}

// BeforeUpdate encryptes data in Database
func (p *Person) BeforeUpdate(_ *gorm.DB) (err error) {
	if p.GivenName == "" || p.LastName == "" {
		return errors.ErrPersonMissingName
	}
	return p.encrypt()
}

// AfterUpdate decryptes data after creation
func (p *Person) AfterUpdate(_ *gorm.DB) (err error) {
	return p.decrypt()
}

// AfterFind decryptes data from Database
func (p *Person) AfterFind(_ *gorm.DB) (err error) {
	return p.decrypt()
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
