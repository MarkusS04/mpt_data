package person

import (
	"encoding/json"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/api/middleware"
	"mpt_data/database/logging"
	"mpt_data/database/person"
	"mpt_data/helper"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"net/http"

	"github.com/gorilla/mux"
)

const packageName = "api.person"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	// person.go
	mux.HandleFunc(apiModel.PersonHref, middleware.CheckAuthentication(getPerson)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.PersonHref, middleware.CheckAuthentication(addPerson)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.PersonHrefWithID, middleware.CheckAuthentication(deletePerson)).Methods(http.MethodDelete)
	mux.HandleFunc(apiModel.PersonHrefWithID, middleware.CheckAuthentication(updatePerson)).Methods(http.MethodPut)

	// task.go
	mux.HandleFunc(apiModel.PersonHrefTask, middleware.CheckAuthentication(getTaskForPerson)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.PersonHrefTask, middleware.CheckAuthentication(addTaskToPerson)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.PersonHrefTask, middleware.CheckAuthentication(deleteTaskFromPerson)).Methods(http.MethodDelete)
}

// @Summary		Get Person
// @Description	Get all Persons
// @Tags			Person
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Success		200	{array}	dbModel.Person
// @Failure		400
// @Failure		401
// @Router			/person [GET]
func getPerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".getPerson"

	tx := middleware.GetTx(r.Context())
	persons, err := person.GetPerson(tx)

	if err != nil {
		api_helper.InternalError(w, funcName, err.Error())
		return
	}
	api_helper.ResponseJSON(w, funcName, persons)
}

// @Summary		Add Person
// @Description	Add Person
// @Tags			Person
// @Accept			json
// @Produce		json
// @Param			Person	body	dbModel.Person	true	"Person"
// @Security		ApiKeyAuth
// @Success		201	{object}	dbModel.Person
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/person [POST]
func addPerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addPerson"
	var personIn dbModel.Person
	if err := json.NewDecoder(r.Body).Decode(&personIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName,
			apiModel.Result{
				Result: "person not created",
				Error:  "provided json data is invalid",
			}, err)
		return
	}

	tx := middleware.GetTx(r.Context())

	err := person.AddPerson(tx, &personIn)
	switch err {
	case nil:
		api_helper.ResponseJSON(w, funcName, personIn, http.StatusCreated)
	case errors.ErrPersonMissingName:
		api_helper.ResponseBadRequest(
			w, funcName, apiModel.Result{
				Result: "failed to store data",
				Error:  err.Error(),
			}, err)
		break
	default:
		api_helper.ResponseBadRequest(
			w, funcName, apiModel.Result{
				Result: "failed to store data",
			}, err)
		break
	}
}

// deletePerson deletes a person
//
//	@Summary		Delete Person
//	@Description	Delete one person with its details
//	@Tags			Person
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"ID of Person"
//	@Security		ApiKeyAuth
//	@Success		200
//	@Failure		400	{object}	apiModel.Result
//	@Failure		401
//	@Router			/person/{id} [DELETE]
func deletePerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deletePerson"

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		api_helper.ResponseBadRequest(w, funcName,
			apiModel.Result{
				Result: "failed to delete person",
				Error:  "id not valid",
			}, err)
		return
	}

	persons := dbModel.Person{}
	persons.ID = uint(*id)

	tx := middleware.GetTx(r.Context())

	if err := person.DeletePerson(tx, persons); err != nil {
		logging.LogError(funcName, err.Error())
		api_helper.ResponseJSON(
			w, funcName,
			apiModel.Result{
				Result: "failed to delete person",
				Error:  "Internal Server Error"})
		return
	}

	api_helper.ResponseJSON(
		w, funcName,
		apiModel.Result{
			Result: "deleted person succesfull"})
}

// updatePerson updates a person
//
//	@Summary		Update Person
//	@Description	Update a person
//	@Tags			Person
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int				true	"ID of Person"
//	@Param			Person	body	dbmodel.Person	true	"Data for Person"
//	@Security		ApiKeyAuth
//	@Success		200
//	@Failure		400	{object}	apiModel.Result
//	@Failure		401
//	@Router			/person/{id} [PUT]
func updatePerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".updatePerson"

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		api_helper.ResponseBadRequest(w, funcName,
			apiModel.Result{
				Result: "person not updated",
				Error:  "id not valid"}, err)
		return
	}

	var personIn dbModel.Person
	if err := json.NewDecoder(r.Body).Decode(&personIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}
	personIn.ID = uint(*id)

	tx := middleware.GetTx(r.Context())

	err = person.UpdatePerson(tx, &personIn)
	switch err {
	case nil:
		api_helper.ResponseJSON(w, funcName, personIn, http.StatusOK)
	case errors.ErrPersonMissingName:
		api_helper.ResponseBadRequest(
			w, funcName, apiModel.Result{
				Result: "failed to store data",
				Error:  err.Error(),
			}, err)
		break
	default:
		api_helper.ResponseBadRequest(
			w, funcName, apiModel.Result{
				Result: "failed to store data",
			}, err)
		break
	}
}
