// Package task provides api routes to manipulate tasks
package task

import (
	"encoding/json"
	"fmt"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/api/middleware"
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/database/task"
	"mpt_data/helper"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"net/http"

	"github.com/gorilla/mux"
)

const packageName = "api.task"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc(apimodel.TaskHref, middleware.CheckAuthentication(getTask)).Methods(http.MethodGet)
	mux.HandleFunc(apimodel.TaskHref, middleware.CheckAuthentication(addTask)).Methods(http.MethodPost)
	mux.HandleFunc(apimodel.TaskHrefWithID, middleware.CheckAuthentication(deleteTask)).Methods(http.MethodDelete)
	mux.HandleFunc(apimodel.TaskHrefWithID, middleware.CheckAuthentication(updateTask)).Methods(http.MethodPut)

	mux.HandleFunc(apimodel.TaskDetailHref, middleware.CheckAuthentication(addTaskDetail)).Methods(http.MethodPost)
	mux.HandleFunc(apimodel.TaskDetailHrefWithID, middleware.CheckAuthentication(deleteTaskDetail)).Methods(http.MethodDelete)
	mux.HandleFunc(apimodel.TaskDetailHrefWithID, middleware.CheckAuthentication(updateTaskDetail)).Methods(http.MethodPut)

	mux.HandleFunc(apimodel.TaskHref, middleware.CheckAuthentication(updateOrderTask)).Methods(http.MethodPut)
	mux.HandleFunc(apimodel.TaskDetailHref, middleware.CheckAuthentication(updateOrderTaskDetail)).Methods(http.MethodPut)
}

// @Summary		Get Task
// @Description	Get all Tasks with their details
// @Tags			Task
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Success		200	{array}		dbModel.Task
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task [GET]
func getTask(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".getTask"

	tx := middleware.GetTx(r.Context())
	tasks, err := task.GetTask(tx)
	if err != nil {
		api_helper.InternalError(w, funcName, err.Error())
		return
	}

	api_helper.ResponseJSON(w, funcName, tasks)
}

// @Summary		Add Task
// @Description	Add Task
// @Tags			Task
// @Accept			json
// @Produce		json
// @Param			Task	body	dbModel.Task	true	"Task"
// @Security		ApiKeyAuth
// @Success		200	{object}	dbModel.Task
// @Success		201
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task [POST]
func addTask(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addTask"
	var taskIn dbModel.Task
	if err := json.NewDecoder(r.Body).Decode(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{
			Error:  "failed to decode request body",
			Result: "Task not added"}, err)
		return
	}

	tx := middleware.GetTx(r.Context())
	err := task.AddTask(tx, &taskIn)
	switch err {
	case nil:
		api_helper.ResponseJSON(w, funcName, apimodel.Result{
			Data:       taskIn,
			Result:     "Task succesfull created",
			StatusCode: 201,
		}, http.StatusCreated)
		break
	case errors.ErrTaskAlreadyExists, errors.ErrTaskDescrNotSet, errors.ErrTaskDetailAlreadyExists:
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{
			Result: "failed to add task",
			Error:  err.Error(),
		}, err)
		break
	default:
		api_helper.InternalError(w, funcName, err.Error())
		break
	}
}

// @Summary		Delete Task
// @Description	Delete one Task with its details
// @Tags			Task
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"ID of task"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task/{id} [DELETE]
func deleteTask(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deleteTask"

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		api_helper.ResponseBadRequest(w, funcName,
			apimodel.Result{
				Result: "failed to delete task",
				Error:  "id not valid"}, err)
		return
	}

	tasks := dbModel.Task{}
	tasks.ID = uint(*id)

	tx := middleware.GetTx(r.Context())
	if err := task.DeleteTask(tx, &tasks); err != nil {
		logging.LogError(funcName, err.Error())
		api_helper.ResponseJSON(w, funcName,
			apimodel.Result{
				Result: "failed to delete task",
				Error:  "Internal Server Error",
			}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Update Task
// @Description	Update the name of a task
// @Tags			Task
// @Accept			json
// @Produce		json
// @Param			id		path	int				true	"ID of task"
// @Param			task	body	dbModel.Task	true	"task"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task/{id} [PUT]
func updateTask(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".updateTask"

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		api_helper.ResponseBadRequest(w, funcName,
			apimodel.Result{
				Result: "task not updated",
				Error:  "id no valid"}, err)
		return
	}

	var taskIn dbModel.Task
	if err := json.NewDecoder(r.Body).Decode(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{
			Result: "task not updated",
			Error:  "failed to decode request body"}, err)
		return
	}
	taskIn.ID = uint(*id)

	tx := middleware.GetTx(r.Context())
	err = task.UpdateTask(tx, taskIn)
	switch err {
	case nil:
		api_helper.ResponseJSON(w, funcName, apimodel.Result{Result: "task sucessful updaed", Data: taskIn}, http.StatusOK)
		break
	case errors.ErrTaskDescrNotSet:
		api_helper.ResponseBadRequest(
			w, funcName, apimodel.Result{
				Result: "task not updated",
				Error:  err.Error(),
			}, err)
		break
	default:
		api_helper.InternalError(w, funcName, fmt.Sprint("failed to update task", err.Error()))
		break
	}
}

// addTaskDetail adds a new detail to a task
//
//	@Summary		Add Detail
//	@Description	Add Detail to Task
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int					true	"ID of task"
//	@Param			TaskDetail	body	dbModel.TaskDetail	true	"TaskDetail"
//	@Security		ApiKeyAuth
//	@Success		201	{object}	dbModel.TaskDetail
//	@Failure		400	{object}	apiModel.Result
//	@Failure		401
//	@Router			/task/{id}/detail [POST]
func addTaskDetail(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addTaskDetail"
	var taskIn dbModel.TaskDetail
	if err := json.NewDecoder(r.Body).Decode(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName,
			apimodel.Result{
				Result: "taskdetail not created",
				Error:  "failed to decode request body"}, err)
		return
	}

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		api_helper.ResponseBadRequest(w, funcName,
			apimodel.Result{
				Result: "taskdetail not created",
				Error:  "id not valid"}, err)
		return
	}

	taskIn.TaskID = uint(*id)

	err = task.AddTaskDetail(middleware.GetTx(r.Context()), &taskIn)
	switch err {
	case nil:
		api_helper.ResponseJSON(w, funcName, apimodel.Result{
			Result:     "Taskdetail succesfull created",
			StatusCode: http.StatusCreated,
		}, http.StatusCreated)
		break
	case errors.ErrTaskDetailAlreadyExists, errors.ErrTaskDescrNotSet:
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{
			Result: "taskdetail not created",
			Error:  err.Error()}, err)
		break
	default:
		api_helper.ResponseJSON(w, funcName, apimodel.Result{
			Result:     "taskdetail not created",
			Error:      "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, http.StatusInternalServerError)
		break
	}
}

// @Summary		Delete Detail
// @Description	Delete Detail of Task
// @Tags			Task
// @Accept			json
// @Produce		json
// @Param			id			path	int	true	"ID of task"
// @Param			detailId	path	int	true	"ID of taskdetail"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task/{id}/detail/{detailId} [DELETE]
func deleteTaskDetail(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deleteTaskDetail"
	var taskIn dbModel.TaskDetail
	idTask, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *idTask <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{
			Result: "failed to delete taskdetail",
			Error:  "id not valid"}, err)
		return
	}

	idDetail, err := helper.ExtractIntFromURL(r, "detailId")
	if err != nil || *idDetail <= 0 {
		api_helper.ResponseBadRequest(w, funcName,
			apimodel.Result{
				Result: "failed to delete taskdetail",
				Error:  "id not valid"}, err)
		return
	}

	taskIn.TaskID = uint(*idTask)
	taskIn.ID = uint(*idDetail)

	err = task.DeleteTaskDetail(middleware.GetTx(r.Context()), taskIn)
	switch err {
	case nil:
		api_helper.ResponseJSON(w, funcName,
			apimodel.Result{
				Result:     "taskdetail succesfull deleted",
				StatusCode: http.StatusOK,
			})
		break
	case errors.ErrTaskDetailAlreadyExists, errors.ErrTaskDescrNotSet:
		api_helper.ResponseJSON(w, funcName,
			apimodel.Result{
				Result:     "failed to delete taskdetail",
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			},
			http.StatusBadRequest)
		break
	default:
		logging.LogError(funcName, err.Error())
		api_helper.ResponseJSON(w, funcName,
			apimodel.Result{
				Result:     "failed to delete taskdetail",
				StatusCode: http.StatusInternalServerError,
			},
			http.StatusInternalServerError)
		break
	}
}

// @Summary		Update TaskDetail
// @Description	Update the name of a taskdetail
// @Tags			Task
// @Accept			json
// @Produce		json
// @Param			id			path	int				true	"ID of task"
// @Param			detailId	path	int				true	"ID of taskdetail"
// @Param			task		body	dbModel.Task	true	"task"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task/{id}/detail/{detailId} [PUT]
func updateTaskDetail(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".updateTask"

	idTask, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *idTask <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "invalid id"}, err)
		return
	}

	idDetail, err := helper.ExtractIntFromURL(r, "detailId")
	if err != nil || *idDetail <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "invalid id"}, err)
		return
	}

	var taskIn dbModel.TaskDetail
	if err := json.NewDecoder(r.Body).Decode(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "failed to decode request body"}, err)
		return
	}
	taskIn.TaskID = uint(*idTask)
	taskIn.ID = uint(*idDetail)

	err = task.UpdateTaskDetail(middleware.GetTx(r.Context()), taskIn)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		break
	case errors.ErrTaskDescrNotSet:
		api_helper.ResponseJSON(w, funcName, apimodel.Result{
			Result:     "failed to update task",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		}, http.StatusBadRequest)
		break
	default:
		api_helper.InternalError(w, funcName, fmt.Sprint("failed to update task", err.Error()))
		break
	}
}

// @Summary		Update Task Order
// @Description	Update the Ordering of Tasks in selects
// @Tags			Task
// @Accept			json
// @Produce		json
// @Param			tasks		body	[]apimodel.OrderTask	true	"Array to hold all tasks and their ordering"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task [PUT]
func updateOrderTask(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".updateOrderTask"

	var data []apimodel.OrderTask
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "failed to decode request body"}, err)
		return
	}

	tx := database.DB.Begin()
	defer tx.Commit()
	if err := task.OrderTask(tx, data); err != nil {
		tx.Rollback()
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "could not update order of task"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Update TaskDetail Order
// @Description	Update the Ordering of TaskDetails in selects
// @Tags			Task
// @Accept			json
// @Produce		json
// @Param			id			path	int				true	"ID of task"
// @Param			tasks		body	[]apimodel.OrderTaskDetail	true	"Array to hold all taskDetails and their ordering"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/task/{id}/detail [PUT]
func updateOrderTaskDetail(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".updateOrderTaskDetail"

	idTask, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *idTask <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "invalid id"}, err)
		return
	}

	var data []apimodel.OrderTaskDetail
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "failed to decode request body"}, err)
		return
	}

	tx := database.DB.Begin()
	defer tx.Commit()
	if err := task.OrderTaskDetail(tx, data, uint(*idTask)); err != nil {
		tx.Rollback()
		api_helper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "could not update order of task"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
