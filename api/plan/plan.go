package plan

import (
	"encoding/json"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/api/auth"
	"mpt_data/database"
	"mpt_data/database/logging"
	"mpt_data/database/plan"
	"mpt_data/helper"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc(apiModel.PlanHref, auth.CheckAuthentication(getPlan)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.PlanHrefWithID, auth.CheckAuthentication(getPlanWithID)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.PlanHrefWithIDPeople, auth.CheckAuthentication(getPersonPlan)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.PlanHref, auth.CheckAuthentication(addPlan)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.PlanHrefWithID, auth.CheckAuthentication(updatePlan)).Methods(http.MethodPut)
}

// @Summary		Get Plan
// @Description	Get Plan for a period
// @Tags			Plan
// @Accept			json
// @Produce		json,application/pdf
// @Param			StartDate	query	string	true	"Start date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Param			EndDate		query	string	true	"End date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Security		ApiKeyAuth
// @Success		200	{array}		dbModel.Plan
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/plan [GET]
func getPlan(w http.ResponseWriter, r *http.Request) {
	funcName := "api.plan.getPlan"
	accept := r.Header.Get("Accept")

	var load func(http.ResponseWriter, *http.Request, time.Time, time.Time)
	// Check if the request wants JSON
	if accept == "application/json" {
		load = getPlanJSON
	} else if accept == "application/pdf" {
		load = getPlanPDF
	} else {
		load = getPlanJSON
		logging.LogWarning(funcName, "accept hedear:"+accept)
	}

	queryParams := r.URL.Query()

	startDate, err := helper.ParseTime(queryParams.Get("StartDate"))
	endDate, err2 := helper.ParseTime(queryParams.Get("EndDate"))
	if err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "could not parse StartDate and/or EndDate"}, err)
		return
	}
	if err2 != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "could not parse StartDate and/or EndDate"}, err2)
		return
	}

	load(w, r, startDate, endDate)
}

func getPlanJSON(w http.ResponseWriter, _ *http.Request, startDate time.Time, endDate time.Time) {
	funcName := "api.plan.getPlanJson"
	plan, err := plan.GetPlan(generalmodel.Period{StartDate: startDate, EndDate: endDate})
	if err != nil {
		api_helper.InternalError(w, funcName, "error creating plan: "+err.Error())
		return
	}
	api_helper.ResponseJSON(w, funcName, plan)
}

// @Summary		Get Plan with ID
// @Description	Get Plan for a specific planId
// @Tags			Plan
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"ID of plan item"
// @Security		ApiKeyAuth
// @Success		200	{object}	dbModel.Plan
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/plan/{id} [GET]
func getPlanWithID(w http.ResponseWriter, r *http.Request) {
	funcName := "api.plan.addPlan"
	if id, err := helper.ExtractIntFromURL(r, "id"); err != nil {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "could not read id in path"}, http.StatusBadRequest)
	} else if *id <= 0 {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "id smaller then 0 invalid"}, http.StatusBadRequest)
	} else {
		if plan, err := plan.GetPlanWithID(uint(*id)); err != nil {
			api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "plan not available"}, http.StatusBadRequest)
		} else {
			api_helper.ResponseJSON(w, funcName, plan)
		}
	}
}

// @Summary		Get people forPlan
// @Description	Loads all people for a meeting with specified task
// @Tags			Plan
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"ID of plan item"
// @Security		ApiKeyAuth
// @Success		200	{array}		dbModel.Person
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/plan/{id}/people [GET]
func getPersonPlan(w http.ResponseWriter, r *http.Request) {
	funcName := "api.plan.addPlan"

	var planData dbModel.Plan
	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "could not read id in path"}, http.StatusBadRequest)
		return
	} else if *id <= 0 {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "id smaller then 0 invalid"}, http.StatusBadRequest)
		return
	}
	if err := database.DB.Preload("Meeting").First(&planData, "id = ?", uint(*id)).Error; err != nil {
		api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "plan not available"}, http.StatusBadRequest)
		return
	}

	people, err := plan.GetAllPersonAvailable(database.DB, planData)
	if err != nil {
		api_helper.InternalError(w, funcName, "error loading persons: "+err.Error())
		return
	}

	api_helper.ResponseJSON(w, funcName, people)
}

// @Summary		Create Plan
// @Description	Create Plan for a period
// @Tags			Plan
// @Accept			json
// @Produce		json
// @Param			StartDate	query	string	true	"Start date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Param			EndDate		query	string	true	"End date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Security		ApiKeyAuth
// @Success		201	{array}		dbModel.Plan
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/plan [POST]
func addPlan(w http.ResponseWriter, r *http.Request) {
	funcName := "api.plan.addPlan"
	queryParams := r.URL.Query()

	startDate, err := helper.ParseTime(queryParams.Get("StartDate"))
	endDate, err2 := helper.ParseTime(queryParams.Get("EndDate"))
	if err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "could not parse StartDate and/or EndDate"}, err)
		return
	}
	if err2 != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "could not parse StartDate and/or EndDate"}, err2)
		return
	}

	plan, err := plan.CreatePlanData(database.DB, generalmodel.Period{StartDate: startDate, EndDate: endDate})
	if err != nil {
		api_helper.InternalError(w, funcName, "error creating plan: "+err.Error())
		return
	}

	api_helper.ResponseJSON(w, funcName, plan)
}

// @Summary		Update a Plan Element
// @Description	Update Person for one task and meeting
// @Tags			Plan
// @Accept			json
// @Produce		json
// @Param			id		path	int					true	"ID of plan item"
// @Param			person	body	plan.updatePlan.p	true	"ID of Person"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400	{object}	apiModel.Result
// @Failure		401
// @Router			/plan/{id} [PUT]
func updatePlan(w http.ResponseWriter, r *http.Request) {
	funcName := "api.plan.updatePlan"

	type p struct {
		ID uint
	}
	var person p
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "error in request body"}, err)
		return
	}

	var planData dbModel.Plan

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "id not correctly set"}, err)
		return
	}
	database.DB.First(&planData, "id = ?", uint(*id))
	planData.PersonID = person.ID

	err = plan.UpdatePlanElement(planData)
	switch err {
	case gorm.ErrRecordNotFound, errors.ErrTaskForPersonNotAllowed:
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		api_helper.InternalError(w, funcName, "error creating plan: "+err.Error())
	}
}
