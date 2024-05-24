package person

import (
	"encoding/json"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/database/person"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"net/http"
)

// @Summary		Get Person-Task
// @Description	Get Tasks of Person(s)
// @Description	ID of Person must always be set, 0 to load all persons with their tasks
// @Tags			Person,Task
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			id	path	int	true	"ID of Person"
// @Success		200 {array} dbModel.Task
// @Failure		400
// @Failure		401
// @Failure		500	"eg. loading failed due to any error"
// @Router			/person/{id}/task [GET]
func getTaskForPerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".getTaskForPerson"
	var (
		// getTaskOnePerson loads all tasks for one Person
		getTaskOnePerson = func(w http.ResponseWriter, id uint) {
			const funcName = packageName + ".getTaskOnePerson"
			tasks, err := person.GetTaskOfPerson(id)
			if err != nil {
				api_helper.InternalError(w, funcName, err.Error())
			}
			api_helper.ResponseJSON(w, funcName, tasks)
		}
		// getTaskAllPerson loads all persons with there tasks
		getTaskAllPerson = func(w http.ResponseWriter) {
			const funcName = packageName + ".getTaskAllPerson"
			persons, err := person.GetPersonWithTask()
			if err != nil {
				api_helper.InternalError(w, funcName, err.Error())
			}
			api_helper.ResponseJSON(w, funcName, persons)
		}
	)

	id, err := api_helper.ExtractIntFromURL(r, "id")
	if err != nil {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: err.Error()}, http.StatusBadRequest)
		return
	}

	if id <= 0 {
		getTaskAllPerson(w)
	} else {
		getTaskOnePerson(w, uint(id))
	}
}

// addTaskToPerson adds new tasks to a person.
//
//	@Summary		Add Tasks to Person
//	@Description	Add Tasks to a Person
//	@Tags			Person,Task
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int		true	"ID of Person"
//	@Param			Task	body	[]uint	true	"Task-Details which should be added to a person"
//	@Security		ApiKeyAuth
//	@Success		201 {array} dbModel.PersonTask
//	@Failure		400
//	@Failure		401
//	@Router			/person/{id}/task [POST]
func addTaskToPerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addTaskToPerson"

	id, err := api_helper.ExtractIntFromURL(r, "id")
	if err != nil {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "invalid ID: " + err.Error()}, http.StatusBadRequest)
		return
	}
	if id <= 0 {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "ID must be greater than 0"}, http.StatusBadRequest)
		return
	}

	var tasks []uint
	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}

	var taskDetails []dbModel.TaskDetail
	for _, task := range tasks {
		taskDetails = append(taskDetails, dbModel.TaskDetail{ID: task})
	}

	addedTasks, err := person.AddTaskToPerson(uint(id), taskDetails)
	if err != nil {
		api_helper.InternalError(w, funcName, "failed to add task to person"+err.Error())
		return
	}

	api_helper.ResponseJSON(w, funcName, addedTasks, http.StatusCreated)
}

// deleteTaskFromPerson deletes tasks of person
//
//	@Summary		Delete Persons Task
//	@Description	Delete tasks of a person
//	@Tags			Person,Task
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int		true	"ID of Person"
//	@Param			Task	body	[]uint	true	"Task-Details which should be deleted from a person"
//	@Security		ApiKeyAuth
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Router			/person/{id}/task [DELETE]
func deleteTaskFromPerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deleteTaskFromPerson"

	id, err := api_helper.ExtractIntFromURL(r, "id")
	if err != nil {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "invalid ID: " + err.Error()}, http.StatusBadRequest)
		return
	}
	if id <= 0 {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "ID must be greater than 0"}, http.StatusBadRequest)
		return
	}

	var tasks []uint
	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}

	var taskDetails []dbModel.TaskDetail
	for _, task := range tasks {
		taskDetails = append(taskDetails, dbModel.TaskDetail{ID: task})
	}

	if err := person.DeleteTaskFromPerson(uint(id), taskDetails); err != nil {
		api_helper.InternalError(w, funcName, "failed delete task from person: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
