package dbmodel

import (
	"encoding/json"
	"mpt_data/helper"

	"time"

	"gorm.io/gorm"
)

type Meeting struct {
	gorm.Model `json:"-"`
	ID         uint
	Date       time.Time `gorm:"uniqueIndex" json:"Date"`
}

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
