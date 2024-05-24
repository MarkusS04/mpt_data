package helper

import "time"

// ParseTime tries to parse the given String with
// "2006-01-02" and "2006-01-02T15:04:05Z07:00"
// returns time if succesfull, otherwise error will be set
func ParseTime(timeString string) (pTime time.Time, err error) {

	layouts := []string{time.DateOnly, time.RFC3339}

	for _, layout := range layouts {
		// Parse the Date field using the desired layout
		date, localerr := time.Parse(layout, timeString)
		if localerr != nil {
			err = localerr
		} else {
			return date, nil
		}
	}
	return time.Time{}, err
}
