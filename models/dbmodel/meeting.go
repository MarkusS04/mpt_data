// Package dbmodel provides all structs for databse ORM
package dbmodel

import (
	"encoding/json"
	"mpt_data/helper"
	"time"

	"gorm.io/gorm"
)

// Meeting stores a date and if set a tag
type Meeting struct {
	gorm.Model `json:"-"`
	ID         uint
	Date       time.Time `gorm:"uniqueIndex" json:"Date"`
	TagID      uint      `json:"-"`
	Tag        Tag       `gorm:"ForeignKey:TagID"`
}

// Tag is a struct to have a descr
type Tag struct {
	gorm.Model `json:"-"`
	ID         uint
	Descr      string
}

// BeforeCreate hook for gorm
func (t *Tag) BeforeCreate(_ *gorm.DB) (err error) {
	tag, err := helper.EncryptData(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = tag
	return
}

// AfterFind hook for gorm
func (t *Tag) AfterFind(_ *gorm.DB) (err error) {
	tag, err := helper.DecryptData(t.Descr)
	if err != nil {
		return err
	}
	t.Descr = string(tag)
	return
}

// UnmarshalJSON unmarshals json as meeting
func (m *Meeting) UnmarshalJSON(data []byte) (err error) {
	// Unmarshal the JSON data into the temporary struct
	var meetingJSON = struct {
		Date string `json:"Date"`
	}{}

	if err := json.Unmarshal(data, &meetingJSON); err != nil {
		return err
	}

	date, err := helper.ParseTime(meetingJSON.Date)
	m.Date = date

	return err
}

// MarshalJSON marshals meeting as json
func (m Meeting) MarshalJSON() ([]byte, error) {
	type Alias Meeting
	if m.Tag.ID == 0 {
		return json.Marshal(&struct {
			Alias
			Tag interface{} `json:",omitempty"`
		}{
			Alias: (Alias)(m),
			Tag:   nil,
		})
	}
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(m),
	})
}
