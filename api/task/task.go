package task

import (
	"encoding/json"
	"fmt"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/api/middleware"
	"mpt_data/database"
	"mpt_data/database/task"
	"mpt_data/helper"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"net/http"

	"github.com/gorilla/mux"
)

const packageName = "api.task"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc(apiModel.TaskHref, middleware.CheckAuthentication(getTask)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.TaskHref, middleware.CheckAuthentication(addTask)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.TaskHrefWithID, middleware.CheckAuthentication(deleteTask)).Methods(http.MethodDelete)
	mux.HandleFunc(apiModel.TaskHrefWithID, middleware.CheckAuthentication(updateTask)).Methods(http.MethodPut)

	mux.HandleFunc(apiModel.TaskDetailHref, middleware.CheckAuthentication(addTaskDetail)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.TaskDetailHrefWithID, middleware.CheckAuthentication(deleteTaskDetail)).Methods(http.MethodDelete)
	mux.HandleFunc(apiModel.TaskDetailHrefWithID, middleware.CheckAuthentication(updateTaskDetail)).Methods(http.MethodPut)

	mux.HandleFunc(apiModel.TaskHref, middleware.CheckAuthentication(updateOrderTask)).Methods(http.MethodPut)
	mux.HandleFunc(apiModel.TaskDetailHref, middleware.CheckAuthentication(updateOrderTaskDetail)).Methods(http.MethodPut)
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
func getTask(w http.ResponseWriter, _ *http.Request) {
	const funcName = packageName + ".getTask"
	tasks, err := task.GetTask()
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
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}

	if err := task.AddTask(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to add task"}, err)
		return
	}

	api_helper.ResponseJSON(w, funcName, taskIn, http.StatusCreated)
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
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	tasks := dbModel.Task{}
	tasks.ID = uint(*id)

	if err := task.DeleteTask(tasks); err != nil {
		api_helper.InternalError(w, funcName, fmt.Sprint("failed to delete task", err.Error()))
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
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	var taskIn dbModel.Task
	if err := json.NewDecoder(r.Body).Decode(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}
	taskIn.ID = uint(*id)

	if err := task.UpdateTask(taskIn); err != nil {
		api_helper.InternalError(w, funcName, fmt.Sprint("failed to update task", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
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
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	taskIn.TaskID = uint(*id)

	if err := task.AddTaskDetail(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to add task"}, err)
		return
	}

	api_helper.ResponseJSON(w, funcName, taskIn, http.StatusCreated)
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
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	idDetail, err := helper.ExtractIntFromURL(r, "detailId")
	if err != nil || *idDetail <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	taskIn.TaskID = uint(*idTask)
	taskIn.ID = uint(*idDetail)

	if err := task.DeleteTaskDetail(taskIn); err != nil {
		api_helper.InternalError(w, funcName, fmt.Sprint("failed to delete task:", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
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
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	idDetail, err := helper.ExtractIntFromURL(r, "detailId")
	if err != nil || *idDetail <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	var taskIn dbModel.TaskDetail
	if err := json.NewDecoder(r.Body).Decode(&taskIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}
	taskIn.TaskID = uint(*idTask)
	taskIn.ID = uint(*idDetail)

	if err := task.UpdateTaskDetail(taskIn); err != nil {
		api_helper.InternalError(w, funcName, fmt.Sprint("failed to update task", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
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

	var data []apiModel.OrderTask
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}

	tx := database.DB.Begin()
	defer tx.Commit()
	if err := task.OrderTask(tx, data); err != nil {
		tx.Rollback()
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "could not update order of task"}, err)
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
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "invalid id"}, err)
		return
	}

	var data []apiModel.OrderTaskDetail
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}

	tx := database.DB.Begin()
	defer tx.Commit()
	if err := task.OrderTaskDetail(tx, data, uint(*idTask)); err != nil {
		tx.Rollback()
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "could not update order of task"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
