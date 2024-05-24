package person

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/database/person"
	"mpt_data/helper"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	api_test "mpt_data/test/api"
	"mpt_data/test/vars"
	"net/http"
	"os"
	"testing"
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

var testPerson = dbModel.Person{GivenName: "Max", LastName: "Mueller"}

func TestAddPerson(t *testing.T) {
	var testcases = []struct {
		name       string
		data       api_test.RequestData
		statusCode int
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   testPerson,
				Route:  apiModel.PersonHref,
				Method: http.MethodPost,
				Router: addPerson,
				Path:   apiModel.PersonHref,
			},
			http.StatusCreated,
		},
	}
	t.Cleanup(func() {
		if person.DeletePerson(testPerson) != nil {
			t.Log("Test cleanup failed")
		}
	})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if status := response.Code; status != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, status)
			}
		})
	}
}

func TestDeletePerson(t *testing.T) {
	if person.AddPerson(&testPerson) != nil {
		t.Errorf("Test preparation failed")
	}
	var testcases = []struct {
		name       string
		data       api_test.RequestData
		statusCode int
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/%d", apiModel.PersonHref, testPerson.ID),
				Method: http.MethodDelete,
				Router: deletePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			http.StatusOK,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if response.Code != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, response.Code)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	if person.AddPerson(&testPerson) != nil {
		t.Errorf("Test preparation failed")
	}
	testPerson.GivenName = "New Name"
	var testcases = []struct {
		name       string
		data       api_test.RequestData
		statusCode int
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   testPerson,
				Route:  fmt.Sprintf("%s/%d", apiModel.PersonHref, testPerson.ID),
				Method: http.MethodPut,
				Router: updatePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			http.StatusOK,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if response.Code != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, response.Code)
			}
		})
	}
}
