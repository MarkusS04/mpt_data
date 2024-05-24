package meeting

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/database/meeting"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	api_test "mpt_data/test/api"
	"mpt_data/test/vars"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if err := database.Connect(vars.GetDbPAth()); err != nil {
		panic(err)
	}
	m.Run()
}

func TestGetMeeting(t *testing.T) {
	// Prepare
	meetings := []dbModel.Meeting{{Date: time.Now()}, {Date: time.Now().AddDate(0, 0, 1)}}
	if meeting.AddMeetings(meetings) != nil {
		t.Errorf("Test preparation failed")
	}

	// testcase
	var testcases = []struct {
		name     string
		data     api_test.RequestData
		meetings []dbModel.Meeting
		err      error
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data: meetings,
				Route: fmt.Sprintf(
					"/api/v1/meeting?StartDate=%s&EndDate=%s",
					time.Now().AddDate(0, 0, -1).Format(time.DateOnly),
					time.Now().AddDate(0, 0, 1).Format(time.DateOnly),
				),
				Method: http.MethodGet,
				Router: getMeetings,
				Path:   apiModel.MeetingHref,
			},
			meetings,
			nil,
		},
	}

	// cleanup
	t.Cleanup(func() {
		database.DB.Unscoped().Delete(&meetings)
	})

	// test
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if status := response.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
				t.Logf("Body: %s", response.Body)
			}
		})
	}
}

func TestAddMeeting(t *testing.T) {
	// Prepare
	meetings := []dbModel.Meeting{{Date: time.Now()}, {Date: time.Now().AddDate(0, 0, 1)}}
	var testcases = []struct {
		name     string
		data     api_test.RequestData
		meetings []dbModel.Meeting
		err      error
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   meetings,
				Route:  "/api/v1/meeting",
				Method: http.MethodPost,
				Router: addMeeting,
				Path:   apiModel.MeetingHref,
			},
			meetings,
			nil,
		},
	}
	t.Cleanup(func() {
		for _, meeting := range meetings {
			database.DB.Unscoped().Delete(&meeting, "date = ?", meeting.Date)
		}
	})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if status := response.Code; status != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, status)
			}
		})
	}
}

func TestUpdateMeeting(t *testing.T) {
	// Prepare
	meetings := dbModel.Meeting{Date: time.Now()}
	if meeting.AddMeetings([]dbModel.Meeting{meetings}) != nil {
		t.Errorf("Test preparation failed")
	}

	meetings.Date = meetings.Date.AddDate(0, 0, 1)
	var testcases = []struct {
		name     string
		data     api_test.RequestData
		meetings dbModel.Meeting
		err      error
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   meetings,
				Route:  "/api/v1/meeting/1",
				Method: http.MethodPut,
				Router: updatetMeeting,
				Path:   apiModel.MeetingHrefWithID,
			},
			meetings,
			nil,
		},
	}
	t.Cleanup(func() {
		database.DB.Unscoped().Delete(&meetings, "date = ?", meetings.Date)
	})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if status := response.Code; status != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, status)
				t.Logf("Body: %s", response.Body)
			}
		})
	}
}

func TestDeleteMeeting(t *testing.T) {
	// Prepare
	var meetings dbModel.Meeting
	if rows := database.DB.First(&meetings).RowsAffected; rows == 0 {
		meetings = dbModel.Meeting{Date: time.Now()}
		if meeting.AddMeetings([]dbModel.Meeting{meetings}) != nil {
			t.Errorf("Test preparation failed")
		}
	}

	var testcases = []struct {
		name    string
		data    api_test.RequestData
		meeting dbModel.Meeting
		err     error
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   meetings,
				Route:  fmt.Sprintf("/api/v1/meeting/%d", meetings.ID),
				Method: http.MethodDelete,
				Router: deleteMeeting,
				Path:   apiModel.MeetingHrefWithID,
			},
			meetings,
			nil,
		},
	}
	t.Cleanup(func() {})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			rr := api_test.DoRequest(t, testcase.data)
			// Assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
				t.Logf("Body: %s", rr.Body)
			}
		})
	}
}
