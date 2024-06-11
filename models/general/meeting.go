// Package generalmodel provides types for db and api
package generalmodel

import "time"

type Period struct {
	StartDate time.Time
	EndDate   time.Time
}
