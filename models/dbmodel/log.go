package dbmodel

import (
	"mpt_data/helper"
	"time"

	"gorm.io/gorm"
)

type Log struct {
	LogLevel  uint
	Source    string
	Text      string
	TimeStamp time.Time
}

// BeforeCreate hook for gorm
func (l *Log) BeforeCreate(_ *gorm.DB) (err error) {
	source, err := helper.EncryptDataToBase64(l.Source)
	if err != nil {
		return err
	}
	l.Source = source
	text, err := helper.EncryptDataToBase64(l.Text)
	if err != nil {
		return err
	}
	l.Text = text
	return
}

// AfterCreate hook for gorm
func (l *Log) AfterCreate(_ *gorm.DB) (err error) {
	source, err := helper.DecryptDataFromBase64(l.Source)
	if err != nil {
		return err
	}
	l.Source = string(source)
	text, err := helper.DecryptDataFromBase64(l.Text)
	if err != nil {
		return err
	}
	l.Text = string(text)
	return
}

// AfterFind hook for gorm
func (l *Log) AfterFind(_ *gorm.DB) (err error) {
	source, err := helper.DecryptDataFromBase64(l.Source)
	if err != nil {
		return err
	}
	l.Source = string(source)
	text, err := helper.DecryptDataFromBase64(l.Text)
	if err != nil {
		return err
	}
	l.Text = string(text)
	return
}
