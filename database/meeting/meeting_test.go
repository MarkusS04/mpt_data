package meeting

import (
	"errors"
	"fmt"
	"mpt_data/database"
	"mpt_data/helper"
	myerrors "mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"mpt_data/test/vars"
	"os"
	"testing"
	"time"
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
	m.Run()
}

func TestGetMeetings(t *testing.T) {
	// Prepare
	meetings := []dbModel.Meeting{
		{Date: time.Now()},
		{Date: time.Now().AddDate(0, 0, 1)},
	}
	database.DB.Create(&meetings)
	t.Cleanup(func() {
		database.DB.Unscoped().Delete(&meetings)
	})
	var testcases = []struct {
		name     string
		period   generalmodel.Period
		meetings []dbModel.Meeting
		err      error
	}{
		{"successful", generalmodel.Period{StartDate: time.Now().AddDate(0, 0, -1), EndDate: time.Now().AddDate(0, 0, 1)}, meetings, nil},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			data, err := GetMeetings(testcase.period)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			if len(data) != len(testcase.meetings) {
				t.Errorf("expected %d, got %d", len(testcase.meetings), len(data))
			}
		})
	}
}

func TestAddMeeting(t *testing.T) {
	// Prepare
	dates := []dbModel.Meeting{{Date: time.Now()}, {Date: time.Now().AddDate(0, 0, 1)}}
	var testcases = []struct {
		name     string
		meetings []dbModel.Meeting
		err      error
	}{
		{"success", dates, nil},
		{"error multiple", dates, errors.New("UNIQUE constraint failed: meetings.id")},
	}

	t.Cleanup(func() {
		database.DB.Unscoped().Delete(dates)
	})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			err := AddMeetings(testcase.meetings)
			// Assert
			if err != testcase.err {
				if err.Error() != testcase.err.Error() {
					t.Errorf("expected %s, got %s", testcase.err, err)
				}
			}
		})
	}
}

func TestUpdateMeeting(t *testing.T) {
	// Prepare
	meeting := dbModel.Meeting{Date: time.Now()}
	if AddMeetings([]dbModel.Meeting{meeting}) != nil {
		t.Errorf("Test preparation failed")
	}
	database.DB.First(&meeting, "date = ?", meeting.Date)

	t.Cleanup(func() {
		database.DB.Unscoped().Delete(&meeting)
	})

	var testcases = []struct {
		name    string
		meeting dbModel.Meeting
		err     error
	}{
		{"successful", meeting, nil},
	}
	t.Cleanup(func() {})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			testcase.meeting.Date = time.Now()
			err := UpdateMeeting(testcase.meeting)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

func TestDeleteMeeting(t *testing.T) {
	// Prepare
	meeting := dbModel.Meeting{Date: time.Now()}
	database.DB.Create(&meeting)

	var testcases = []struct {
		name    string
		meeting dbModel.Meeting
		err     error
	}{
		{"succes", meeting, nil},
		{"error", meeting, myerrors.ErrMeetingNotDeleted},
	}
	t.Cleanup(func() {})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			err := DeleteMeeting(testcase.meeting)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}
