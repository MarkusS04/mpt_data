// Package dbmodel provides all structs for databse ORM
package dbmodel

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Meeting stores a date and if set a tag
type Meeting struct {
	gorm.Model `json:"-"`
	ID         uint
	Date       time.Time `gorm:"uniqueIndex" json:"Date"`
	TagID      uint      `json:"-"`
	Tag        Tag       `gorm:"ForeignKey:TagID" json:"Tag,omitempty"`
}

// Tag is a struct to have a descr
type Tag struct {
	gorm.Model `json:"-"`
	ID         uint
	Descr      string
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
