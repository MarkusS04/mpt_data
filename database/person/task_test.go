package person

import (
	"mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	"testing"
)

func TestAddTaskToPerson(t *testing.T) {
	var testcases = []struct {
		name       string
		personID   uint
		taskDetail []dbModel.TaskDetail
		err        error
	}{
		{"succesfull", 1, []dbModel.TaskDetail{{ID: 1}}, nil},
		{"error person not set", 0, []dbModel.TaskDetail{{ID: 1}}, errors.ErrIDNotSet},
		{"error person not set", 1, []dbModel.TaskDetail{{ID: 0}}, errors.ErrIDNotSet},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			data, err := AddTaskToPerson(testcase.personID, testcase.taskDetail)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
				return
			}
			if err == nil && data[0].ID == 0 {
				t.Errorf("expected data to be set")
			}
		})
	}
}

func TestDeleteTaskFromPerson(t *testing.T) {
	var testcases = []struct {
		name       string
		personID   uint
		taskDetail []dbModel.TaskDetail
		err        error
	}{
		{"succesfull", 1, []dbModel.TaskDetail{{ID: 1}}, nil},
		{"error person not set", 0, []dbModel.TaskDetail{{ID: 1}}, errors.ErrIDNotSet},
		{"error person not set", 1, []dbModel.TaskDetail{{ID: 0}}, errors.ErrIDNotSet},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			err := DeleteTaskFromPerson(testcase.personID, testcase.taskDetail)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}
