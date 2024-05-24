package task

import (
	"errors"
	"fmt"
	"mpt_data/database"
	"mpt_data/helper"
	myerrors "mpt_data/helper/errors"
	dbModel "mpt_data/models/dbmodel"
	database_test "mpt_data/test/database"
	"mpt_data/test/vars"
	"os"
	"testing"

	"gorm.io/gorm"
)

var testTask = &dbModel.Task{
	Descr: "Test",
	TaskDetails: []dbModel.TaskDetail{
		{Descr: "Test A"},
		{Descr: "Test B"},
	},
}

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

func TestAddTask(t *testing.T) {
	var testcases = []struct {
		name              string
		model             dbModel.Task
		err               error
		dbCountTask       int64
		dbCountTaskDetail int64
	}{
		{"successfull", *testTask, nil, 1, 2},
		{"error duplicate", *testTask, errors.New("UNIQUE constraint failed: tasks.descr"), 0, 0},
	}
	t.Cleanup(func() {
		if DeleteTask(*testTask) != nil {
			t.Log("Test cleanup failed")
		}
	})

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			countBeforeT := database_test.CountEntries(&dbModel.Task{})
			countBeforeTD := database_test.CountEntries(&dbModel.TaskDetail{})
			// Act
			err := AddTask(&testcase.model)
			countAfterT := database_test.CountEntries(&dbModel.Task{})
			countAfterTD := database_test.CountEntries(&dbModel.TaskDetail{})
			// Assert
			if err != testcase.err {
				if testcase.err == nil {
					t.Errorf("expected nil, got %s", err)
				} else if testcase.err.Error() != err.Error() {
					t.Errorf("expected %s, got %s", testcase.err, err)
				}
			}

			if countAfterT-countBeforeT != testcase.dbCountTask {
				t.Errorf("expected %d tasks, got %d", testcase.dbCountTask, countAfterT-countBeforeT)
			}
			if countAfterTD-countBeforeTD != testcase.dbCountTaskDetail {
				t.Errorf("expected %d taskdetails, got %d", testcase.dbCountTaskDetail, countAfterTD-countBeforeTD)
				t.Logf("before: %d", countBeforeTD)
				t.Logf("after: %d", countAfterTD)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	// Prepare
	countBeforeT := database_test.CountEntries(&dbModel.Task{})
	countBeforeTD := database_test.CountEntries(&dbModel.TaskDetail{})
	if err := AddTask(testTask); err != nil {
		t.Skip("no testdata available")
	}
	var testcases = []struct {
		name  string
		model dbModel.Task
		err   error
	}{
		{"successfull", *testTask, nil},
		{"not available", *testTask, nil},
		{"not available", dbModel.Task{Descr: "bla"}, gorm.ErrRecordNotFound},
	}

	// testing
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			err := DeleteTask(testcase.model)
			countAfterT := database_test.CountEntries(&dbModel.Task{})
			countAfterTD := database_test.CountEntries(&dbModel.TaskDetail{})
			// Assert
			if err != testcase.err {
				if testcase.err == nil {
					t.Errorf("expected nil, got %s", err)
				} else if err == nil {
					t.Errorf("expected %s, got nil", testcase.err)
				} else if testcase.err.Error() != err.Error() {
					t.Errorf("expected %s, got %s", testcase.err, err)
				}
			}
			if countBeforeT != countAfterT {
				t.Errorf("expected task to be deleted, but didn't")
			}
			if countBeforeTD != countAfterTD {
				t.Errorf("expected taskDetail to be deleted, but didn't")
			}
		})
	}
}

func TestAddTaskDetail(t *testing.T) {
	if AddTask(testTask) != nil {
		t.Errorf("Test preparation failed")
	}
	taskDetail := dbModel.TaskDetail{Descr: "Test C", TaskID: testTask.ID}
	var testcases = []struct {
		name  string
		task  dbModel.TaskDetail
		err   error
		count int64
	}{
		{"succesfull", taskDetail, nil, 1},
	}
	t.Cleanup(func() {
		if DeleteTask(*testTask) != nil {
			t.Log("Test cleanup failed")
		}
	})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			countBefore := database_test.CountEntries(&dbModel.TaskDetail{})
			err := AddTaskDetail(&taskDetail)
			countAfter := database_test.CountEntries(&dbModel.TaskDetail{})
			// Assert
			if err != testcase.err {
				t.Errorf("expexted %s, got %s", testcase.err, err)
			}
			if countAfter-countBefore != testcase.count {
				t.Errorf("expected %d entries created, got %d", testcase.count, countAfter-countBefore)
			}
		})
	}

	t.Run("error id null", func(t *testing.T) {
		err := AddTaskDetail(&dbModel.TaskDetail{Descr: "Test"})
		if err != myerrors.ErrForeignIDNotSet {
			t.Errorf("expected %s, got %s", myerrors.ErrForeignIDNotSet, err)
		}
	})
}
func TestDeleteTaskDetail(t *testing.T) {
	if AddTask(testTask) != nil {
		t.Errorf("Test preparation failed")
	}
	var testcases = []struct {
		name  string
		task  dbModel.TaskDetail
		err   error
		count int64
	}{
		{"succesfull", testTask.TaskDetails[0], nil, 1},
	}
	t.Cleanup(func() {
		if DeleteTask(*testTask) != nil {
			t.Log("Test cleanup failed")
		}
	})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			countBefore := database_test.CountEntries(&dbModel.TaskDetail{})
			err := DeleteTaskDetail(testcase.task)
			countAfter := database_test.CountEntries(&dbModel.TaskDetail{})
			// Assert
			if err != testcase.err {
				t.Errorf("expexted %s, got %s", testcase.err, err)
			}
			if countBefore-countAfter != testcase.count {
				t.Errorf("expected %d entries created, got %d", testcase.count, countBefore-countAfter)
			}
		})
	}

	t.Run("error id null", func(t *testing.T) {
		err := AddTaskDetail(&dbModel.TaskDetail{Descr: "Test"})
		if err != myerrors.ErrForeignIDNotSet {
			t.Errorf("expected %s, got %s", myerrors.ErrForeignIDNotSet, err)
		}
	})
}
func TestUpdateTask(t *testing.T) {
	// prepare
	var testcases = []struct {
		name string
		task dbModel.Task
		id   uint
		err  error
	}{
		{"sucessfull", *testTask, 1, nil},
		{"no id set", *testTask, 0, myerrors.ErrIDNotSet},
	}
	if AddTask(testTask) != nil {
		t.Errorf("Test preparation failed")
	}
	testTaskID := testTask.ID
	t.Cleanup(func() {
		testTask.ID = testTaskID
		if DeleteTask(*testTask) != nil {
			t.Log("Test cleanup failed")
		}
	})

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			testcase.task.Descr = "Test2"
			testcase.task.ID = testcase.id

			err := UpdateTask(testcase.task)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

func TestUpdateTaskDetail(t *testing.T) {
	// prepare
	var testcases = []struct {
		name string
		task dbModel.TaskDetail
		id   uint
		err  error
	}{
		{"sucessfull", testTask.TaskDetails[0], 1, nil},
		{"no id set", testTask.TaskDetails[0], 0, myerrors.ErrIDNotSet},
	}
	if AddTask(testTask) != nil {
		t.Errorf("Test preparation failed")
	}
	testTaskID := testTask.ID
	t.Cleanup(func() {
		testTask.ID = testTaskID
		if DeleteTask(*testTask) != nil {
			t.Log("Failed cleanup test")
		}
	})

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			testcase.task.Descr = "Test2"
			testcase.task.ID = testcase.id

			err := UpdateTaskDetail(testcase.task)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}
