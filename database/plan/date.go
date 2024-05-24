package plan

import (
	"time"
)

// languageCode
const (
	German = "de_DE"
)

func getMonthName(month time.Month, lang string) string {
	// Map für die deutschen Wochentage
	monthNames := map[string](map[time.Month]string){
		German: {
			time.January:   "Januar",
			time.February:  "Februar",
			time.March:     "März",
			time.April:     "April",
			time.May:       "Mai",
			time.June:      "Juni",
			time.July:      "Juli",
			time.August:    "August",
			time.September: "September",
			time.October:   "Oktober",
			time.November:  "November",
			time.December:  "Dezember",
		}}

	// Den Namen des Wochentags aus der Map abrufen
	return monthNames[lang][month]
}

func getWeekdayName(weekday time.Weekday, lang string) string {
	// Map für die deutschen Wochentage
	dayName := map[string](map[time.Weekday]string){
		German: {
			time.Monday:    "Montag",
			time.Tuesday:   "Dienstag",
			time.Wednesday: "Mittwoch",
			time.Thursday:  "Donnerstag",
			time.Friday:    "Freitag",
			time.Saturday:  "Samstag",
			time.Sunday:    "Sonntag",
		}}

	// Den Namen des Wochentags aus der Map abrufen
	return dayName[lang][weekday]
}
