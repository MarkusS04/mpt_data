package person

import (
	"encoding/json"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/api/auth"
	"mpt_data/database"
	"mpt_data/database/person"
	"mpt_data/helper"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"net/http"

	"github.com/gorilla/mux"
)

const packageName = "api.person"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	// person.go
	mux.HandleFunc(apiModel.PersonHref, auth.CheckAuthentication(getPerson)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.PersonHref, auth.CheckAuthentication(addPerson)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.PersonHrefWithID, auth.CheckAuthentication(deletePerson)).Methods(http.MethodDelete)
	mux.HandleFunc(apiModel.PersonHrefWithID, auth.CheckAuthentication(updatePerson)).Methods(http.MethodPut)

	// task.go
	mux.HandleFunc(apiModel.PersonHrefTask, auth.CheckAuthentication(getTaskForPerson)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.PersonHrefTask, auth.CheckAuthentication(addTaskToPerson)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.PersonHrefTask, auth.CheckAuthentication(deleteTaskFromPerson)).Methods(http.MethodDelete)
}

//	@Summary		Get Person
//	@Description	Get all Persons
//	@Tags			Person
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{array}	dbModel.Person
//	@Failure		400
//	@Failure		401
//	@Router			/person [GET]
func getPerson(w http.ResponseWriter, _ *http.Request) {
	const funcName = packageName + ".getPerson"
	persons, err := person.GetPerson(database.DB)
	if err != nil {
		api_helper.InternalError(w, funcName, err.Error())
		return
	}
	api_helper.ResponseJSON(w, funcName, persons)
}

//	@Summary		Add Person
//	@Description	Add Person
//	@Tags			Person
//	@Accept			json
//	@Produce		json
//	@Param			Person	body	dbModel.Person	true	"Person"
//	@Security		ApiKeyAuth
//	@Success		201	{object}	dbModel.Person
//	@Failure		400	{object}	apiModel.Result
//	@Failure		401
//	@Router			/person [POST]
func addPerson(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addPerson"
	var personIn dbModel.Person
	if err := json.NewDecoder(r.Body).Decode(&personIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	if err := person.AddPerson(tx, &personIn); err != nil {
		tx.Rollback()
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to store data"}, err)
		return
	}

	api_helper.ResponseJSON(w, funcName, personIn, http.StatusCreated)
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
	if err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to extract id"}, err)
		return
	}

	if *id <= 0 {
		http.Error(w, "invalid id in URL", http.StatusBadRequest)
		return
	}

	persons := dbModel.Person{}
	persons.ID = uint(*id)

	tx := database.DB.Begin()
	defer tx.Commit()

	if err := person.DeletePerson(tx, persons); err != nil {
		tx.Rollback()
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to delete person"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	if err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to extract id"}, err)
		return
	}

	if *id <= 0 {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "invalid id in URL"}, http.StatusBadRequest)
		return
	}

	var personIn dbModel.Person
	if err := json.NewDecoder(r.Body).Decode(&personIn); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "failed to decode request body"}, err)
		return
	}
	personIn.ID = uint(*id)

	tx := database.DB.Begin()
	defer tx.Commit()

	if err := person.UpdatePerson(tx, &personIn); err != nil {
		tx.Rollback()
		api_helper.InternalError(w, funcName, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
