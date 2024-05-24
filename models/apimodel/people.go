// Package apimodel contains all models used by api to exchange data with client
package apimodel

import "mpt_data/models/dbmodel"

// People is a struct to hold absent and available people
type People struct {
	Absent    []dbmodel.Person `json:"absent"`
	Available []dbmodel.Person `json:"available"`
	Assigned  dbmodel.Person   `json:"assigned"`
}
