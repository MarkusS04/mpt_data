package task

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/database/task"
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

var (
	testTask = dbModel.Task{
		Descr: "Test Task",
		TaskDetails: []dbModel.TaskDetail{
			{Descr: "Detail 1"},
			{Descr: "Detail 2"},
		},
	}
	testTaskDetail = dbModel.TaskDetail{Descr: "Hello World"}
)

func TestAddTask(t *testing.T) {
	var testcases = []struct {
		name       string
		data       api_test.RequestData
		statusCode int
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   testTask,
				Route:  apiModel.TaskHref,
				Method: http.MethodPost,
				Router: addTask,
				Path:   apiModel.TaskHref,
			},
			http.StatusCreated,
		},
	}
	t.Cleanup(func() {
		if task.DeleteTask(testTask) != nil {
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

func TestDeleteTask(t *testing.T) {
	if task.AddTask(&testTask) != nil {
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
				Route:  fmt.Sprintf("%s/%d", apiModel.TaskHref, testTask.ID),
				Method: http.MethodDelete,
				Router: deleteTask,
				Path:   apiModel.TaskHrefWithID,
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
	t.Cleanup(func() {
		if task.DeleteTask(testTask) != nil {
			t.Log("Test cleanup failed")
		}
	})

	if task.AddTask(&testTask) != nil {
		t.Errorf("Test preparation failed")
	}
	testTask.Descr = "New Text"
	var testcases = []struct {
		name       string
		data       api_test.RequestData
		statusCode int
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   testTask,
				Route:  fmt.Sprintf("%s/%d", apiModel.TaskHref, testTask.ID),
				Method: http.MethodPut,
				Router: updateTask,
				Path:   apiModel.TaskHrefWithID,
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

func TestAddTaskDetail(t *testing.T) {
	if task.AddTask(&testTask) != nil {
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
				Data:   testTaskDetail,
				Route:  fmt.Sprintf("%s/%d/detail", apiModel.TaskHref, testTask.ID),
				Method: http.MethodPost,
				Router: addTaskDetail,
				Path:   apiModel.TaskDetailHref,
			},
			http.StatusCreated,
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

func TestDeleteTaskDetail(t *testing.T) {
	t.Cleanup(func() {
		if task.DeleteTask(testTask) != nil {
			t.Log("Test cleanup failed")
		}
	})
	var testcases = []struct {
		name       string
		data       api_test.RequestData
		statusCode int
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/%d/detail/%d", apiModel.TaskHref, testTask.ID, testTask.TaskDetails[0].ID),
				Method: http.MethodDelete,
				Router: deleteTaskDetail,
				Path:   apiModel.TaskDetailHrefWithID,
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
func TestUpdateTaskDetail(t *testing.T) {
	if task.AddTask(&testTask) != nil {
		t.Errorf("Test preparation failed")
	}
	testTask.TaskDetails[0].Descr = "New Text"
	var testcases = []struct {
		name       string
		data       api_test.RequestData
		statusCode int
	}{
		{
			"succesfull",
			api_test.RequestData{
				Data:   testTask,
				Route:  fmt.Sprintf("%s/%d/detail/%d", apiModel.TaskHref, testTask.ID, testTask.TaskDetails[0].ID),
				Method: http.MethodPut,
				Router: updateTaskDetail,
				Path:   apiModel.TaskDetailHrefWithID,
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
