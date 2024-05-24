package dbmodel

import "time"

type Log struct {
	LogLevel  uint
	Source    string
	Text      string
	TimeStamp time.Time
}
