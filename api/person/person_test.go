package person

import (
	"encoding/json"
	"fmt"
	"mpt_data/database"
	"mpt_data/database/person"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	apiModel "mpt_data/models/apimodel"
	"mpt_data/models/dbmodel"
	dbModel "mpt_data/models/dbmodel"
	api_test "mpt_data/test/api"
	"mpt_data/test/vars"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	vars.PrepareConfig()
	m.Run()
}

func TestAddPerson(t *testing.T) {
	var testcases = []struct {
		name         string
		data         api_test.RequestData
		statusCode   int
		errorMessage string
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   dbModel.Person{GivenName: "Max", LastName: "Mueller"},
				Route:  apiModel.PersonHref,
				Method: http.MethodPost,
				Router: addPerson,
				Path:   apiModel.PersonHref,
			},
			http.StatusCreated,
			"",
		},
		{
			"no givenname",
			api_test.RequestData{
				Data:   dbModel.Person{GivenName: "", LastName: "Mueller"},
				Route:  apiModel.PersonHref,
				Method: http.MethodPost,
				Router: addPerson,
				Path:   apiModel.PersonHref,
			},
			http.StatusBadRequest,
			errors.ErrPersonMissingName.Error(),
		},
		{
			"no lastname",
			api_test.RequestData{
				Data:   dbModel.Person{GivenName: "Max", LastName: ""},
				Route:  apiModel.PersonHref,
				Method: http.MethodPost,
				Router: addPerson,
				Path:   apiModel.PersonHref,
			},
			http.StatusBadRequest,
			errors.ErrPersonMissingName.Error(),
		},
		{
			"no name",
			api_test.RequestData{
				Data:   dbModel.Person{GivenName: "", LastName: ""},
				Route:  apiModel.PersonHref,
				Method: http.MethodPost,
				Router: addPerson,
				Path:   apiModel.PersonHref,
			},
			http.StatusBadRequest,
			errors.ErrPersonMissingName.Error(),
		},
		{
			"invalid json",
			api_test.RequestData{
				Data:   "\"data\": key",
				Route:  apiModel.PersonHref,
				Method: http.MethodPost,
				Router: addPerson,
				Path:   apiModel.PersonHref,
			},
			http.StatusBadRequest,
			"provided json data is invalid",
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if status := response.Code; status != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, status)
			}
			if testcase.errorMessage != "" {
				var result apimodel.Result
				if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
					t.Errorf("decoding result failed: %v", err)
				}

				if result.Error != testcase.errorMessage {
					t.Errorf("expected error %s, got %s", testcase.errorMessage, result.Error)
				}
			}
		})
	}
}

func TestDeletePerson(t *testing.T) {
	var testcases = []struct {
		name          string
		person        dbmodel.Person
		data          api_test.RequestData
		statusCode    int
		errorMessage  string
		resultMessage string
		addPerson     bool
	}{
		{
			name: "succesfull",
			person: dbModel.Person{
				GivenName: "Tester",
				LastName:  "Test",
			},
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/", apiModel.PersonHref),
				Method: http.MethodDelete,
				Router: deletePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode:    http.StatusOK,
			resultMessage: "deleted person succesfull",
			addPerson:     true,
		},
		{
			name: "invalid id",
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/bla", apiModel.PersonHref),
				Method: http.MethodDelete,
				Router: deletePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode:    http.StatusBadRequest,
			addPerson:     false,
			errorMessage:  "id not valid",
			resultMessage: "failed to delete person",
		},
		{
			name: "no id",
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/", apiModel.PersonHref),
				Method: http.MethodDelete,
				Router: deletePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode: http.StatusNotFound,
			addPerson:  false,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.addPerson {
				person.AddPerson(database.DB, &testcase.person)
				testcase.data.Route += fmt.Sprint(testcase.person.ID)
				defer person.DeletePerson(database.DB, testcase.person)
			}
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if response.Code != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, response.Code)
			}

			var result apimodel.Result
			if testcase.errorMessage != "" || testcase.resultMessage != "" {
				if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
					t.Errorf("decoding result failed: %v", err)
				}
			}
			if testcase.errorMessage != "" {
				if result.Error != testcase.errorMessage {
					t.Errorf("expected error %s, got %s", testcase.errorMessage, result.Error)
				}
			}
			if testcase.resultMessage != "" {
				if result.Result != testcase.resultMessage {
					t.Errorf("expected error %s, got %s", testcase.resultMessage, result.Result)
				}
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	var testcases = []struct {
		name          string
		person        dbmodel.Person
		data          api_test.RequestData
		statusCode    int
		errorMessage  string
		resultMessage string
		addPerson     bool
	}{
		{
			name: "succesfull",
			person: dbModel.Person{
				GivenName: "Tester",
				LastName:  "Test",
			},
			data: api_test.RequestData{
				Data:   dbModel.Person{GivenName: "Max", LastName: "Mueller"},
				Route:  apiModel.PersonHref + "/",
				Method: http.MethodPut,
				Router: updatePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode: http.StatusOK,
			addPerson:  true,
		},
		{
			name: "no givenname",
			data: api_test.RequestData{
				Data:   dbModel.Person{GivenName: "", LastName: "Mueller"},
				Route:  apiModel.PersonHref + "/",
				Method: http.MethodPut,
				Router: updatePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode:   http.StatusBadRequest,
			errorMessage: errors.ErrPersonMissingName.Error(),
			person: dbModel.Person{
				GivenName: "Tester",
				LastName:  "Test",
			},
			addPerson: true,
		},
		{
			name: "no lastname",
			data: api_test.RequestData{
				Data:   dbModel.Person{GivenName: "Max", LastName: ""},
				Route:  apiModel.PersonHref + "/",
				Method: http.MethodPut,
				Router: updatePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode:   http.StatusBadRequest,
			errorMessage: errors.ErrPersonMissingName.Error(),
			person: dbModel.Person{
				GivenName: "Tester",
				LastName:  "Test",
			},
			addPerson: true,
		},
		{
			name: "no name",
			data: api_test.RequestData{
				Data:   dbModel.Person{GivenName: "", LastName: ""},
				Route:  apiModel.PersonHref + "/",
				Method: http.MethodPut,
				Router: updatePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode:   http.StatusBadRequest,
			errorMessage: errors.ErrPersonMissingName.Error(),
			person: dbModel.Person{
				GivenName: "Tester",
				LastName:  "Test",
			},
			addPerson: true,
		},
		{
			name: "no id with invalid data",
			data: api_test.RequestData{
				Data:   "\"data\": key",
				Route:  apiModel.PersonHref + "/",
				Method: http.MethodPut,
				Router: updatePerson,
				Path:   apiModel.PersonHrefWithID,
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.addPerson {
				person.AddPerson(database.DB, &testcase.person)
				testcase.data.Route += fmt.Sprint(testcase.person.ID)
				defer person.DeletePerson(database.DB, testcase.person)
			}
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if response.Code != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, response.Code)
			}

			var result apimodel.Result
			if testcase.errorMessage != "" || testcase.resultMessage != "" {
				if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
					t.Errorf("decoding result failed: %v", err)
				}
			}
			if testcase.errorMessage != "" {
				if result.Error != testcase.errorMessage {
					t.Errorf("expected error %s, got %s", testcase.errorMessage, result.Error)
				}
			}
			if testcase.resultMessage != "" {
				if result.Result != testcase.resultMessage {
					t.Errorf("expected error %s, got %s", testcase.resultMessage, result.Result)
				}
			}
		})
	}
}
