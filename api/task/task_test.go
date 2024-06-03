package task

import (
	"encoding/json"
	"fmt"
	"mpt_data/database"
	"mpt_data/database/task"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	apiModel "mpt_data/models/apimodel"
	"mpt_data/models/dbmodel"
	api_test "mpt_data/test/api"
	"mpt_data/test/vars"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	vars.PrepareConfig()
	m.Run()
}

var testTaskDetail = dbmodel.TaskDetail{Descr: "Hello World"}

func TestAddTask(t *testing.T) {
	var testcases = []struct {
		name          string
		data          api_test.RequestData
		errorMessage  string
		resultMessage string
		statusCode    int
	}{
		{
			name: "succesfull",
			data: api_test.RequestData{
				Data: dbmodel.Task{
					Descr: "Test Task",
					TaskDetails: []dbmodel.TaskDetail{
						{Descr: "Detail 1"},
						{Descr: "Detail 2"},
					},
				},
				Route:  apiModel.TaskHref,
				Method: http.MethodPost,
				Router: addTask,
				Path:   apiModel.TaskHref,
			},
			statusCode:    http.StatusCreated,
			errorMessage:  "",
			resultMessage: "Task succesfull created",
		},
		{
			name: "descr missing",
			data: api_test.RequestData{
				Data: dbmodel.Task{
					Descr: "",
					TaskDetails: []dbmodel.TaskDetail{
						{Descr: "Detail 1"},
					},
				},
				Route:  apiModel.TaskHref,
				Method: http.MethodPost,
				Router: addTask,
				Path:   apiModel.TaskHref,
			},
			statusCode:    http.StatusBadRequest,
			errorMessage:  "task or taskdetail descr missing",
			resultMessage: "failed to add task",
		},
		{
			name: "descr missing",
			data: api_test.RequestData{
				Data: dbmodel.Task{
					Descr: "Task",
					TaskDetails: []dbmodel.TaskDetail{
						{Descr: ""},
					},
				},
				Route:  apiModel.TaskHref,
				Method: http.MethodPost,
				Router: addTask,
				Path:   apiModel.TaskHref,
			},
			statusCode:    http.StatusBadRequest,
			errorMessage:  "task or taskdetail descr missing",
			resultMessage: "failed to add task",
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
			if testcase.errorMessage != "" || testcase.resultMessage != "" {
				var result apimodel.Result
				if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
					t.Errorf("decoding result failed: %v", err)
				}

				if result.Error != testcase.errorMessage {
					t.Errorf("expected error %s, got %s", testcase.errorMessage, result.Error)
				}
				if result.Result != testcase.resultMessage {
					t.Errorf("expected result %s, got %s", testcase.resultMessage, result.Result)
				}
			}
		})
	}

	t.Run("duplicate", func(t *testing.T) {
		// Prepare
		tasks := &dbmodel.Task{
			Descr: "Test",
			TaskDetails: []dbmodel.TaskDetail{
				{Descr: "Detail 1"},
			},
		}
		task.AddTask(database.DB, tasks)
		defer task.DeleteTask(database.DB, tasks)
		// Act
		response := api_test.DoRequest(t, api_test.RequestData{
			Data: dbmodel.Task{
				Descr: "Test",
				TaskDetails: []dbmodel.TaskDetail{
					{Descr: "Detail 1"},
				},
			},
			Route:  apiModel.TaskHref,
			Method: http.MethodPost,
			Router: addTask,
			Path:   apiModel.TaskHref,
		})
		// Assert
		if status := response.Code; status != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, status)
		}

		var result apimodel.Result
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			t.Errorf("decoding result failed: %v", err)
		}

		if result.Error != "task with descr already exists" {
			t.Errorf("expected error %s, got %s", "task with descr already exists", result.Error)
		}
		if result.Result != "failed to add task" {
			t.Errorf("expected result %s, got %s", "failed to add task", result.Result)
		}

	})
}

func TestDeleteTask(t *testing.T) {
	var testcases = []struct {
		name          string
		task          dbmodel.Task
		data          api_test.RequestData
		statusCode    int
		errorMessage  string
		resultMessage string
		addTask       bool
	}{
		{
			name: "succesfull",
			task: dbmodel.Task{
				Descr: "Task",
				TaskDetails: []dbmodel.TaskDetail{
					{Descr: "Detail 1"},
				},
			},
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/", apiModel.TaskHref),
				Method: http.MethodDelete,
				Router: deleteTask,
				Path:   apiModel.TaskHrefWithID,
			},
			statusCode:    http.StatusOK,
			resultMessage: "deleted task succesfull",
			addTask:       true,
		},
		{
			name: "invalid id",
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/bla", apiModel.TaskHref),
				Method: http.MethodDelete,
				Router: deleteTask,
				Path:   apiModel.TaskHrefWithID,
			},
			statusCode:    http.StatusBadRequest,
			resultMessage: "failed to delete task",
			errorMessage:  "id not valid",
		},
		{
			name: "no id",
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/", apiModel.TaskHref),
				Method: http.MethodDelete,
				Router: deleteTask,
				Path:   apiModel.TaskHrefWithID,
			},
			statusCode: http.StatusNotFound,
			addTask:    false,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.addTask {
				task.AddTask(database.DB, &testcase.task)
				testcase.data.Route += fmt.Sprint(testcase.task.ID)
				defer task.DeleteTask(database.DB, &testcase.task)
			}
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if response.Code != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, response.Code)
			}

			var result apimodel.Result
			if testcase.errorMessage == "" || testcase.resultMessage == "" {
				return
			}
			if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
				t.Errorf("decoding result failed: %v", err)
			}

			if result.Error != testcase.errorMessage {
				t.Errorf("expected error %s, got %s", testcase.errorMessage, result.Error)
			}

			if result.Result != testcase.resultMessage {
				t.Errorf("expected error %s, got %s", testcase.resultMessage, result.Result)
			}

		})
	}
}

func TestUpdateTask(t *testing.T) {
	var testcases = []struct {
		name          string
		task          dbmodel.Task
		data          api_test.RequestData
		statusCode    int
		errorMessage  string
		resultMessage string
		addTask       bool
	}{
		{
			name: "succesfull",
			task: dbmodel.Task{
				Descr: "Test",
				TaskDetails: []dbmodel.TaskDetail{
					{Descr: "Test"},
				},
			},
			data: api_test.RequestData{
				Data:   dbmodel.Task{Descr: "Bla"},
				Route:  apiModel.TaskHref + "/",
				Method: http.MethodPut,
				Router: updateTask,
				Path:   apiModel.TaskHrefWithID,
			},
			statusCode: http.StatusOK,
			addTask:    true,
		},
		{
			name: "no descr",
			task: dbmodel.Task{
				Descr: "Test",
				TaskDetails: []dbmodel.TaskDetail{
					{Descr: "Test"},
				},
			},
			data: api_test.RequestData{
				Data:   dbmodel.Task{Descr: ""},
				Route:  apiModel.TaskHref + "/",
				Method: http.MethodPut,
				Router: updateTask,
				Path:   apiModel.TaskHrefWithID,
			},
			statusCode:   http.StatusBadRequest,
			errorMessage: errors.ErrTaskDescrNotSet.Error(),
			addTask:      true,
		},
		{
			name: "no id with invalid data",
			data: api_test.RequestData{
				Data:   "\"data\": key",
				Route:  apiModel.TaskHref + "/",
				Method: http.MethodPut,
				Router: updateTask,
				Path:   apiModel.TaskHrefWithID,
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.addTask {
				task.AddTask(database.DB, &testcase.task)
				testcase.data.Route += fmt.Sprint(testcase.task.ID)
				defer task.DeleteTask(database.DB, &testcase.task)
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

func TestAddTaskdetail(t *testing.T) {
	var testcases = []struct {
		name          string
		task          *dbmodel.Task
		data          api_test.RequestData
		errorMessage  string
		resultMessage string
		statusCode    int
		addRoute      bool
	}{
		{
			name: "succesfull",
			task: &dbmodel.Task{Descr: "Task"},
			data: api_test.RequestData{
				Data:   dbmodel.TaskDetail{Descr: "Detail 1"},
				Route:  apiModel.TaskHref,
				Method: http.MethodPost,
				Router: addTaskDetail,
				Path:   apiModel.TaskDetailHref,
			},
			statusCode:    http.StatusCreated,
			errorMessage:  "",
			resultMessage: "Taskdetail succesfull created",
			addRoute:      true,
		},
		{
			name: "descr missing",
			task: &dbmodel.Task{Descr: "Task"},
			data: api_test.RequestData{
				Data:   dbmodel.TaskDetail{Descr: ""},
				Route:  apiModel.TaskHref,
				Method: http.MethodPost,
				Router: addTask,
				Path:   apiModel.TaskDetailHref,
			},
			statusCode:    http.StatusBadRequest,
			errorMessage:  "task or taskdetail descr missing",
			resultMessage: "failed to add task",
			addRoute:      true,
		},
		{
			name: "task id invalid",
			task: &dbmodel.Task{Descr: "Task"},
			data: api_test.RequestData{
				Data:   dbmodel.TaskDetail{Descr: "Bla"},
				Route:  apiModel.TaskHref + "/0/detail",
				Method: http.MethodPost,
				Router: addTaskDetail,
				Path:   apiModel.TaskDetailHref,
			},
			statusCode:    http.StatusBadRequest,
			errorMessage:  "id not valid",
			resultMessage: "taskdetail not created",
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			if err := task.AddTask(database.DB, testcase.task); err != nil {
				t.Skipf("test prep failed: %v", err)
			}
			if testcase.addRoute {
				testcase.data.Route += fmt.Sprintf("/%v/detail", testcase.task.ID)
			}
			defer task.DeleteTask(database.DB, testcase.task)
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if status := response.Code; status != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, status)
			}
			if testcase.errorMessage != "" || testcase.resultMessage != "" {
				var result apimodel.Result
				if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
					t.Errorf("decoding result failed: %v", err)
				}

				if result.Error != testcase.errorMessage {
					t.Errorf("expected error %s, got %s", testcase.errorMessage, result.Error)
				}
				if result.Result != testcase.resultMessage {
					t.Errorf("expected result %s, got %s", testcase.resultMessage, result.Result)
				}
			}
		})
	}

	t.Run("duplicate", func(t *testing.T) {
		// Prepare
		tasks := &dbmodel.Task{
			Descr: "Test",
			TaskDetails: []dbmodel.TaskDetail{
				{Descr: "Detail 1"},
			},
		}
		task.AddTask(database.DB, tasks)
		defer task.DeleteTask(database.DB, tasks)
		// Act
		response := api_test.DoRequest(t, api_test.RequestData{
			Data:   dbmodel.TaskDetail{Descr: "Detail 1"},
			Route:  fmt.Sprintf("%s/%d/detail", apimodel.TaskHref, tasks.ID),
			Method: http.MethodPost,
			Router: addTaskDetail,
			Path:   apiModel.TaskDetailHref,
		})
		// Assert
		if status := response.Code; status != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, status)
		}

		var result apimodel.Result
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			t.Errorf("decoding result failed: %v", err)
		}

		if result.Error != "taskdetail with descr already exists" {
			t.Errorf("expected error %s, got %s", "taskdetail with descr already exists", result.Error)
		}
		if result.Result != "taskdetail not created" {
			t.Errorf("expected result %s, got %s", "taskdetail not created", result.Result)
		}

	})
}

func TestDeleteTaskdetail(t *testing.T) {
	var testcases = []struct {
		name          string
		task          *dbmodel.Task
		data          api_test.RequestData
		statusCode    int
		errorMessage  string
		resultMessage string
		addTask       bool
	}{
		{
			name: "succesfull",
			task: &dbmodel.Task{
				Descr: "Task",
				TaskDetails: []dbmodel.TaskDetail{
					{Descr: "Detail 1"},
				},
			},
			data: api_test.RequestData{
				Route:  apiModel.TaskHref,
				Method: http.MethodDelete,
				Router: deleteTaskDetail,
				Path:   apiModel.TaskDetailHrefWithID,
			},
			statusCode:    http.StatusOK,
			resultMessage: "deleted taskdetail succesfull",
			addTask:       true,
		},
		{
			name: "invalid id",
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/10/detail/bla", apiModel.TaskHref),
				Method: http.MethodDelete,
				Router: deleteTaskDetail,
				Path:   apiModel.TaskDetailHrefWithID,
			},
			statusCode:    http.StatusBadRequest,
			resultMessage: "failed to delete taskdetail",
			errorMessage:  "id not valid",
		},
		{
			name: "no id",
			data: api_test.RequestData{
				Data:   nil,
				Route:  fmt.Sprintf("%s/10/detail/", apiModel.TaskHref),
				Method: http.MethodDelete,
				Router: deleteTaskDetail,
				Path:   apiModel.TaskHrefWithID,
			},
			statusCode: http.StatusNotFound,
			addTask:    false,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			var taskdetail dbmodel.TaskDetail
			if testcase.addTask {
				task.AddTask(database.DB, testcase.task)
				taskdetail = testcase.task.TaskDetails[0]
				testcase.data.Route += fmt.Sprintf("/%d/detail/%d", testcase.task.ID, taskdetail.ID)
				defer task.DeleteTask(database.DB, testcase.task)
			}
			// Act
			response := api_test.DoRequest(t, testcase.data)
			// Assert
			if response.Code != testcase.statusCode {
				t.Errorf("expected %d, got %d", testcase.statusCode, response.Code)
			}

			var result apimodel.Result
			if testcase.errorMessage == "" || testcase.resultMessage == "" {
				return
			}
			if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
				t.Errorf("decoding result failed: %v", err)
			}

			if result.Error != testcase.errorMessage {
				t.Errorf("expected error %s, got %s", testcase.errorMessage, result.Error)
			}

			if result.Result != testcase.resultMessage {
				t.Errorf("expected error %s, got %s", testcase.resultMessage, result.Result)
			}

		})
	}
}

func TestUpdateTaskdetail(t *testing.T) {
	var testcases = []struct {
		name          string
		task          dbmodel.Task
		data          api_test.RequestData
		statusCode    int
		errorMessage  string
		resultMessage string
		addTask       bool
	}{
		{
			name: "succesfull",
			task: dbmodel.Task{
				Descr: "Test",
				TaskDetails: []dbmodel.TaskDetail{
					{Descr: "Test"},
				},
			},
			data: api_test.RequestData{
				Data:   dbmodel.TaskDetail{Descr: "Bla"},
				Route:  apiModel.TaskHref + "/",
				Method: http.MethodPut,
				Router: updateTaskDetail,
				Path:   apiModel.TaskDetailHrefWithID,
			},
			statusCode: http.StatusOK,
			addTask:    true,
		},
		{
			name: "no descr",
			task: dbmodel.Task{
				Descr: "Test",
				TaskDetails: []dbmodel.TaskDetail{
					{Descr: "Test"},
				},
			},
			data: api_test.RequestData{
				Data:   dbmodel.Task{Descr: ""},
				Route:  apiModel.TaskHref + "/",
				Method: http.MethodPut,
				Router: updateTaskDetail,
				Path:   apiModel.TaskDetailHrefWithID,
			},
			statusCode:   http.StatusBadRequest,
			errorMessage: errors.ErrTaskDescrNotSet.Error(),
			addTask:      true,
		},
		{
			name: "no id with invalid data",
			data: api_test.RequestData{
				Data:   "\"data\": key",
				Route:  apiModel.TaskHref + "/10/detail/",
				Method: http.MethodPut,
				Router: updateTaskDetail,
				Path:   apiModel.TaskDetailHrefWithID,
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.addTask {
				task.AddTask(database.DB, &testcase.task)
				testcase.data.Route += fmt.Sprintf("%d/detail/%d", testcase.task.ID, testcase.task.TaskDetails[0].ID)
				defer task.DeleteTask(database.DB, &testcase.task)
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
