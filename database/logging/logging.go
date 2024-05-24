package logging

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/helper"
	dbModel "mpt_data/models/dbmodel"
	"os"
	"time"
)

type LogLevel uint

const (
	Info = iota
	Warning
	Error
)

// logtype: info, warning, error
// source: which package and funtion created the log
// text: text to be logged
// TimeStamp: when did the error occur
// Entry: string that will be logged in file
type log struct {
	LogLevel  LogLevel
	Source    string
	Text      string
	TimeStamp time.Time
	Entry     string
}

func addLog(LogLevel LogLevel, Source, Text string) {
	log := log{
		TimeStamp: time.Now(),
		LogLevel:  LogLevel,
		Source:    Source,
		Text:      Text,
	}

	year, month, day := log.TimeStamp.Date()

	logFilename := fmt.Sprintf("%s/%d_%d_%d.log", helper.Config.Log.Path, year, month, day)
	file, err := os.OpenFile(logFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	log.Entry = fmt.Sprintf("%s;%d;%s;%s;\n", log.TimeStamp, LogLevel, log.Source, log.Text)
	if _, err := file.WriteString(log.Entry); err != nil {
		fmt.Println("Error:", err)
	}

	if helper.Config.Log.LevelDB <= uint(LogLevel) {
		dbLog(log)
	}
}

func LogInfo(Source, Text string) {
	addLog(Info, Source, Text)
}
func LogWarning(Source, Text string) {
	addLog(Warning, Source, Text)
}
func LogError(Source, Text string) {
	addLog(Error, Source, Text)
}

func dbLog(log log) {
	db := database.DB.Begin()
	defer db.Commit()
	if err :=
		db.Create(&dbModel.Log{
			LogLevel:  uint(log.LogLevel),
			Source:    log.Source,
			Text:      log.Text,
			TimeStamp: log.TimeStamp,
		}).Error; err != nil {
		fmt.Println("could not create db log, Error: ", err)
		db.Rollback()
	}
}
