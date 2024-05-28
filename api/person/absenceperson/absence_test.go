package absenceperson

import (
	"encoding/json"
	"fmt"
	"mpt_data/database"
	"mpt_data/helper/config"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	api_test "mpt_data/test/api"
	database_test "mpt_data/test/database"
	"mpt_data/test/vars"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	meetingT = &dbModel.Meeting{
		Date: time.Now(), ID: 1,
	}
	absenceT = dbModel.PersonAbsence{
		PersonID: 1,
	}
)

func TestMain(m *testing.M) {
	if err := database.Connect(vars.GetDbPAth()); err != nil {
		panic(err)
	}
	if err := os.Chdir("../../.."); err != nil {
		fmt.Println(err)
	}
	// Load the config
	config.LoadConfig()

	if err := database.DB.Create(&meetingT).Error; err != nil {
		fmt.Println("err", err)
	}
	absenceT.MeetingID = meetingT.ID
	m.Run()
	database.DB.Unscoped().Delete(&meetingT)
}

func TestAddAbsence(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name string
		data api_test.RequestData
		err  error
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   []uint{1},
				Route:  fmt.Sprintf("/api/v1/person/%d/absence", absenceT.PersonID),
				Method: http.MethodPost,
				Router: addAbsence,
				Path:   apiModel.PersonAbsence,
			},
			nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			countBefore := database_test.CountEntries(&dbModel.PersonAbsence{})
			response := api_test.DoRequest(t, testcase.data)
			countAfter := database_test.CountEntries(&dbModel.PersonAbsence{})
			// Assert
			if status := response.Code; status != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, status)
				t.Log(response.Body)
			}
			if countAfter-countBefore != 1 {
				t.Errorf("expected 1 entry to be created, got %d", countAfter-countBefore)
			}
		})
	}
}

func TestGetAbsences(t *testing.T) {
	// Prepare
	var meeting dbModel.PersonAbsence
	database.DB.
		Preload("Person").
		Preload("Meeting").
		Where("meeting_id = ? and person_id=?", meetingT.ID, absenceT.PersonID).
		First(&meeting)

	ab := meeting.Meeting
	// testcase
	var testcases = []struct {
		name     string
		data     api_test.RequestData
		response []dbModel.Meeting
		err      error
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data: nil,
				Route: fmt.Sprintf("/api/v1/person/%d/absence?StartDate=%s&EndDate=%s",
					absenceT.PersonID,
					time.Now().AddDate(0, 0, -1).Format(time.DateOnly),
					time.Now().AddDate(0, 0, 1).Format(time.DateOnly),
				),
				Method: http.MethodGet,
				Router: getAbsence,
				Path:   apiModel.PersonAbsence,
			},
			[]dbModel.Meeting{ab},
			nil,
		},
	}

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
			respBody := response.Body.String()
			if json, _ := json.Marshal(testcase.response); string(json) != respBody {
				t.Errorf("expected %s, got %s", string(json), respBody)
			}
		})
	}
}

func TestDeleteAbsence(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name string
		data api_test.RequestData
		err  error
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   []uint{1},
				Route:  fmt.Sprintf("/api/v1/person/%d/absence", absenceT.PersonID),
				Method: http.MethodDelete,
				Router: deleteAbsence,
				Path:   apiModel.PersonAbsence,
			},
			nil,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			countBefore := database_test.CountEntries(&dbModel.PersonAbsence{})
			rr := api_test.DoRequest(t, testcase.data)
			countAfter := database_test.CountEntries(&dbModel.PersonAbsence{})
			// Assert

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
				t.Logf("Body: %s", rr.Body)
			}
			if countBefore-countAfter != 1 {
				t.Errorf("expected 1 entry to be created, got %d", countBefore-countAfter)
			}
		})
	}
}
