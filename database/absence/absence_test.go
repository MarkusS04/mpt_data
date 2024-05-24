package absence

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/helper"
	dbModel "mpt_data/models/dbmodel"
	database_test "mpt_data/test/database"
	"mpt_data/test/vars"
	"os"
	"testing"
	"time"
)

var (
	meetingT = &dbModel.Meeting{
		Date: time.Now(), ID: 1,
	}
	absence = dbModel.PersonAbsence{
		PersonID: 1,
	}
)

func TestMain(m *testing.M) {
	if err := database.Connect(vars.GetDbPAth()); err != nil {
		panic(err)
	}
	if err := os.Chdir("../.."); err != nil {
		fmt.Println(err)
	}
	// Load the config
	helper.LoadConfig()

	if err := database.DB.Create(&meetingT).Error; err != nil {
		fmt.Println("err", err)
	}
	absence.MeetingID = meetingT.ID
	m.Run()
	database.DB.Unscoped().Delete(&meetingT)
}

func TestAddAbsence(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name    string
		absence []dbModel.PersonAbsence
		err     error
		nums    int64
	}{
		{"success", []dbModel.PersonAbsence{absence}, nil, 1},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			countBefore := database_test.CountEntries(&dbModel.PersonAbsence{})

			err := AddAbsence(testcase.absence)
			countAfter := database_test.CountEntries(&dbModel.PersonAbsence{})
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			if countAfter-countBefore != testcase.nums {
				t.Errorf("expected %d entries created, got %d", testcase.nums, countAfter-countBefore)
			}
		})
	}
}

func TestGetAbsences(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name    string
		absence []dbModel.PersonAbsence
		err     error
	}{
		{"successful", []dbModel.PersonAbsence{absence}, nil},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			data, err := GetAbsenceMeeting(meetingT.ID)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			if len(data) != len(testcase.absence) {
				t.Errorf("expected %d, got %d", len(testcase.absence), len(data))
			}
		})
	}
}

func TestDeleteMeeting(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name    string
		absence []dbModel.PersonAbsence
		err     error
	}{
		{"succes", []dbModel.PersonAbsence{absence}, nil},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			err := DeleteAbsence(testcase.absence)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}
