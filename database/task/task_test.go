package task

import (
	"mpt_data/database"
	myerrors "mpt_data/helper/errors"
	"mpt_data/models/dbmodel"
	dbModel "mpt_data/models/dbmodel"
	"mpt_data/test/vars"
	"testing"
)

func TestMain(m *testing.M) {
	vars.PrepareConfig()
	m.Run()
}

func TestAddTask(t *testing.T) {
	var testcases = []struct {
		name              string
		task              *dbModel.Task
		err               error
		dbCountTask       int64
		dbCountTaskDetail int64
	}{
		{
			name: "successfull",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
					{Descr: "Test B"},
				},
			},
			dbCountTask:       1,
			dbCountTaskDetail: 2,
		},
		{
			name: "Taskdetail-Desct not set",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: ""},
				},
			},
			err: myerrors.ErrTaskDescrNotSet,
		},
		{
			name: "Task-Descr not set",
			task: &dbModel.Task{
				Descr: "",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
					{Descr: "Test B"},
				},
			},
			err: myerrors.ErrTaskDescrNotSet,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()
			// Act
			err := AddTask(tx, testcase.task)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %v, got %v", testcase.err, err)
			}

			if testcase.dbCountTask != 0 || testcase.dbCountTaskDetail != 0 {
				var count int64
				tx.Find(&dbmodel.Task{}).Count(&count)
				if count != testcase.dbCountTask {
					t.Errorf("expected %d, got %d tasks", testcase.dbCountTask, count)
				}

				tx.Find(&dbmodel.TaskDetail{}).Count(&count)
				if count != testcase.dbCountTaskDetail {
					t.Errorf("expected %d, got %d tasks", testcase.dbCountTaskDetail, count)
				}
			}
		})
	}

	t.Run("duplicate task", func(t *testing.T) {
		// Prepare
		tx := database.DB.Begin()
		defer tx.Rollback()
		task := &dbModel.Task{
			Descr: "Test",
			TaskDetails: []dbModel.TaskDetail{
				{Descr: "Test A"},
				{Descr: "Test B"},
			},
		}
		err := AddTask(tx, task)
		if err != nil {
			t.Skipf("test preparation failed: %v", err)
		}
		task.ID = 0

		// Act
		err = AddTask(tx, task)
		// Assert
		if err != myerrors.ErrTaskAlreadyExists {
			t.Errorf("expected %v, got %v", myerrors.ErrTaskAlreadyExists, err)
		}

	})
}
func TestUpdateTask(t *testing.T) {
	var testcases = []struct {
		name     string
		task     *dbModel.Task
		newDescr string
		setID    bool
		err      error
	}{
		{
			name: "successfull",
			task: &dbModel.Task{
				Descr: "Test",
			},
			newDescr: "Test 2",
			setID:    true,
		},
		{
			name: "id not set",
			task: &dbModel.Task{
				Descr: "Test",
			},
			newDescr: "Test 2",
			err:      myerrors.ErrIDNotSet,
		},
		{
			name: "Task-Descr not set",
			task: &dbModel.Task{
				Descr: "Test",
			},
			newDescr: "",
			setID:    true,
			err:      myerrors.ErrTaskDescrNotSet,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()

			err := AddTask(tx, testcase.task)
			if err != nil {
				t.Skipf("failed to prepare test: %v", err)
			}
			// Act
			if !testcase.setID {
				testcase.task.ID = 0
			}

			testcase.task.Descr = testcase.newDescr
			err = UpdateTask(tx, *testcase.task)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}
func TestDeleteTask(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name  string
		task  *dbModel.Task
		setID bool
		err   error
	}{
		{
			name: "successfull",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
					{Descr: "Test B"},
				},
			},
			setID: true,
		},
		{
			name: "id not set",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
					{Descr: "Test B"},
				},
			},
			err: myerrors.ErrIDNotSet,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()
			AddTask(tx, testcase.task)
			// Act
			if !testcase.setID {
				testcase.task.ID = 0
			}
			err := DeleteTask(tx, testcase.task)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

func TestAddTaskDetail(t *testing.T) {
	var testcases = []struct {
		name              string
		task              *dbModel.Task
		taskdetail        *dbModel.TaskDetail
		err               error
		dbCountTask       int64
		dbCountTaskDetail int64
	}{
		{
			name: "successfull",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
				},
			},
			taskdetail: &dbModel.TaskDetail{
				Descr: "Test B",
			},
			dbCountTaskDetail: 2,
		},
		{
			name: "Taskdetail-Descr not set",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: ""},
				},
			},
			err: myerrors.ErrTaskDescrNotSet,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()
			err := AddTask(tx, testcase.task)
			if err != nil {
				t.Skipf("Test prep failed: %v", err)
			}
			// Act
			testcase.taskdetail.TaskID = testcase.task.ID
			err = AddTaskDetail(tx, testcase.taskdetail)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %v, got %v", testcase.err, err)
			}

			if testcase.dbCountTaskDetail != 0 {
				var count int64
				tx.Find(&dbmodel.TaskDetail{}).Count(&count)
				if count != testcase.dbCountTaskDetail {
					t.Errorf("expected %d, got %d tasks", testcase.dbCountTaskDetail, count)
				}
			}
		})
	}

	t.Run("duplicate task detail", func(t *testing.T) {
		// Prepare
		tx := database.DB.Begin()
		defer tx.Rollback()
		task := &dbModel.Task{
			Descr: "Test",
			TaskDetails: []dbModel.TaskDetail{
				{Descr: "Test A"},
			},
		}
		err := AddTask(tx, task)
		if err != nil {
			t.Skipf("test preparation failed: %v", err)
		}

		// Act
		err = AddTaskDetail(tx, &dbModel.TaskDetail{TaskID: task.ID, Descr: "Test A"})
		// Assert
		if err != myerrors.ErrTaskDetailAlreadyExists {
			t.Errorf("expected %v, got %v", myerrors.ErrTaskDetailAlreadyExists, err)
		}
	})

	t.Run("id not set", func(t *testing.T) {
		// Prepare
		tx := database.DB.Begin()
		defer tx.Rollback()
		task := &dbModel.Task{
			Descr: "Test",
			TaskDetails: []dbModel.TaskDetail{
				{Descr: "Test A"},
			},
		}
		err := AddTask(tx, task)
		if err != nil {
			t.Skipf("test preparation failed: %v", err)
		}

		// Act
		err = AddTaskDetail(tx, &dbModel.TaskDetail{TaskID: 0, Descr: "Test A"})
		// Assert
		if err != myerrors.ErrForeignIDNotSet {
			t.Errorf("expected %v, got %v", myerrors.ErrForeignIDNotSet, err)
		}
	})
}
func TestUpdateTaskDetail(t *testing.T) {
	var testcases = []struct {
		name     string
		task     *dbModel.Task
		newDescr string
		setID    bool
		err      error
	}{
		{
			name: "successfull",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
				},
			},
			newDescr: "Testing Tasks",
			setID:    true,
		},
		{
			name: "Taskdetail-Descr not set",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Descr"},
				},
			},
			setID: true,
			err:   myerrors.ErrTaskDescrNotSet,
		},
		{
			name: "id not set",
			task: &dbModel.Task{
				Descr: "",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
				},
			},
			err: myerrors.ErrForeignIDNotSet,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()

			if err := AddTask(tx, testcase.task); err != nil {
				t.Skipf("test preperation failed: %v", err)
			}
			// Act
			taskDetail := testcase.task.TaskDetails[0]
			taskDetail.Descr = testcase.newDescr

			if !testcase.setID {
				taskDetail.ID = 0
			}

			err := UpdateTaskDetail(tx, taskDetail)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}
func TestDeleteTaskDetail(t *testing.T) {
	var testcases = []struct {
		name       string
		task       *dbModel.Task
		setID      bool
		countAfter int64
		err        error
	}{
		{
			name: "successfull",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
				},
			},
			countAfter: 0,
			setID:      true,
		},
		{
			name: "id not set",
			task: &dbModel.Task{
				Descr: "Test",
				TaskDetails: []dbModel.TaskDetail{
					{Descr: "Test A"},
				},
			},
			countAfter: 1,
			err:        myerrors.ErrIDNotSet,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()
			AddTask(tx, testcase.task)
			// Act
			taskDetail := testcase.task.TaskDetails[0]
			if !testcase.setID {
				taskDetail.ID = 0
			}
			err := DeleteTaskDetail(tx, taskDetail)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}

			var count int64
			tx.Find(&dbmodel.TaskDetail{}).Count(&count)
			if count != testcase.countAfter {
				t.Errorf("expexted %d, got %d taskdetails to exits", testcase.countAfter, count)
			}
		})
	}
}
