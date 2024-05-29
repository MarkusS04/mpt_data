package absenceperson

import (
	"encoding/json"
	"mpt_data/api/apihelper"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/api/middleware"
	"mpt_data/database"
	"mpt_data/database/absence"
	"mpt_data/helper"
	"mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

const packageName = "api.person.absence"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc(apimodel.PersonAbsence, middleware.CheckAuthentication(getAbsence)).Methods(http.MethodGet)
	mux.HandleFunc(apimodel.PersonAbsence, middleware.CheckAuthentication(addAbsence)).Methods(http.MethodPost)
	mux.HandleFunc(apimodel.PersonAbsence, middleware.CheckAuthentication(deleteAbsence)).Methods(http.MethodDelete)

	mux.HandleFunc(apimodel.PersonAbsenceRecuring, middleware.CheckAuthentication(getAbsenceRecurring)).Methods(http.MethodGet)
	mux.HandleFunc(apimodel.PersonAbsenceRecuring, middleware.CheckAuthentication(addAbsenceRecurring)).Methods(http.MethodPost)
	mux.HandleFunc(apimodel.PersonAbsenceRecuring, middleware.CheckAuthentication(deleteAbsenceRecurring)).Methods(http.MethodDelete)
}

// @Summary		Get Absence
// @Description	Get absence of person in period
// @Tags			Person,Absence
// @Accept			json
// @Produce		json
// @Param			PersonId	path	int		true	"ID of person"
// @Param			StartDate	query	string	true	"Start date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Param			EndDate		query	string	true	"End date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Security		ApiKeyAuth
// @Success		200	{array}	dbModel.Meeting
// @Failure		400
// @Failure		401
// @Router			/person/{PersonId}/absence [GET]
func getAbsence(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".getAbsence"
	var id *int
	var err error
	if id, err = helper.ExtractIntFromURL(r, "id"); err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "id not valid"}, err)
		return
	}

	queryParams := r.URL.Query()

	startDate, err := helper.ParseTime(queryParams.Get("StartDate"))
	endDate, err2 := helper.ParseTime(queryParams.Get("EndDate"))
	if err != nil || err2 != nil {
		apihelper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "could not parse StartDate and/or EndDate"}, err)
		return
	}

	data, err := absence.GetAbsencePerson(uint(*id), generalmodel.Period{StartDate: startDate, EndDate: endDate})
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusOK)
		case gorm.ErrEmptySlice, gorm.ErrInvalidData:
			w.WriteHeader(http.StatusBadRequest)
		default:
			api_helper.InternalError(w, funcName, err.Error())
		}
	} else {
		api_helper.ResponseJSON(w, funcName, data)
	}
}

// @Summary		Add Absence
// @Description	Add absence of person in period
// @Tags			Person,Absence
// @Accept			json
// @Produce		json
// @Param			PersonId	path	int		true	"ID of person"
// @Param			Absence		body	[]uint	true	"ID of meetings where person is absent"
// @Security		ApiKeyAuth
// @Success		200
// @Success		201
// @Failure		400
// @Failure		401
// @Router			/person/{PersonId}/absence [POST]
func addAbsence(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addAbsence"
	var absences []uint
	var absencePerson []dbModel.PersonAbsence
	if err := json.NewDecoder(r.Body).Decode(&absences); err != nil {
		apihelper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "meetings id not valid"}, err)
		return
	}

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "id not valid"}, err)
		return
	}
	for _, ab := range absences {
		absencePerson = append(absencePerson,
			dbModel.PersonAbsence{PersonID: uint(*id), MeetingID: ab})
	}

	err = absence.AddAbsence(absencePerson)
	if err != nil {
		switch err {
		case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusBadRequest)
		default:
			api_helper.InternalError(w, funcName, err.Error())
		}
	} else {
		api_helper.ResponseJSON(w, funcName, absences, http.StatusCreated)
	}
}

// @Summary		Delete Absence
// @Description	Delete absence of person in period
// @Tags			Person,Absence
// @Accept			json
// @Produce		json
// @Param			PersonId	path	int		true	"ID of person"
// @Param			Absence		body	[]uint	true	"ID of meeting where person is no longer absent"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400
// @Failure		401
// @Router			/person/{PersonId}/absence [DELETE]
func deleteAbsence(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deleteAbsence"
	var absences []uint
	var absencePerson []dbModel.PersonAbsence
	if err := json.NewDecoder(r.Body).Decode(&absences); err != nil {
		apihelper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "meetings id not valid"}, err)
		return
	}

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "people id not valid"}, err)
		return
	}
	for _, ab := range absences {
		absencePerson = append(absencePerson,
			dbModel.PersonAbsence{PersonID: uint(*id), MeetingID: ab})
	}

	err = absence.DeleteAbsence(absencePerson)
	if err != nil {
		switch err {
		case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusBadRequest)
		default:
			api_helper.InternalError(w, funcName, err.Error())
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary		Get Absence
// @Description	Get recurring absence of person
// @Tags			Person,Absence
// @Accept			json
// @Produce		json
// @Param			PersonId	path	int	true	"ID of person"
// @Security		ApiKeyAuth
// @Success		200	{array}	dbModel.PersonRecurringAbsence
// @Failure		400
// @Failure		401
// @Router			/person/{PersonId}/absencerecurring [GET]
func getAbsenceRecurring(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".getAbsenceRecurring"
	var id *int
	var err error
	if id, err = helper.ExtractIntFromURL(r, "id"); err != nil || *id <= 0 {
		apihelper.ResponseError(w, funcName, *apimodel.GetInavalidRequestProblemDetails(
			http.StatusBadRequest,
			"",
			"",
			err,
		))
		return
	}

	data, err := absence.GetRecurringAbsence(uint(*id), database.DB)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusOK)
		case gorm.ErrEmptySlice, gorm.ErrInvalidData:
			w.WriteHeader(http.StatusBadRequest)
		default:
			api_helper.InternalError(w, funcName, err.Error())
		}
	} else {
		api_helper.ResponseJSON(w, funcName, data)
	}
}

// @Summary		Add Absence
// @Description	Add recurring absence to person
// @Tags			Person,Absence
// @Accept			json
// @Produce		json
// @Param			PersonId	path	int		true	"ID of person"
// @Param			Absence		body	[]int	true	"Weekdays where person is absent. 0 = Sunday"
// @Security		ApiKeyAuth
// @Success		200
// @Success		201
// @Failure		400
// @Failure		401
// @Router			/person/{PersonId}/absencerecurring [POST]
func addAbsenceRecurring(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addAbsenceRecurring"

	var absences []int
	if err := json.NewDecoder(r.Body).Decode(&absences); err != nil {
		api_helper.ResponseError(w, funcName, *apimodel.GetInavalidRequestProblemDetails(
			http.StatusBadRequest,
			"Weekdays not correctly specified",
			"", nil))
		return
	}

	var absencePerson []*dbModel.PersonRecurringAbsence
	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		api_helper.ResponseError(w, funcName, *apimodel.GetInavalidRequestProblemDetails(
			http.StatusBadRequest,
			"ID not correctly specified",
			"", nil,
		))
		return
	}
	for _, ab := range absences {
		if ab < 0 || ab > 6 {
			continue
		}
		absencePerson = append(absencePerson,
			&dbModel.PersonRecurringAbsence{PersonID: uint(*id), Weekday: ab})
	}

	db := database.DB.Begin()
	defer db.Commit()

	err = absence.AddRecurringAbsence(absencePerson, db)
	if err != nil {
		db.Rollback()
		switch err {
		case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusBadRequest)
		default:
			api_helper.InternalError(w, funcName, err.Error())
		}
	} else {
		api_helper.ResponseJSON(w, funcName, absences, http.StatusCreated)
	}
}

// @Summary		Delete Absence
// @Description	Delete absence of person in period
// @Tags			Person,Absence
// @Accept			json
// @Produce		json
// @Param			PersonId	path	int		true	"ID of person"
// @Param			Absence		body	[]uint	true	"Weekday where person is no longer absent"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400
// @Failure		401
// @Router			/person/{PersonId}/absencerecurring [DELETE]
func deleteAbsenceRecurring(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deleteAbsenceRecurring"

	idPerson, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *idPerson <= 0 {
		api_helper.ResponseError(w, funcName, *apimodel.GetInavalidRequestProblemDetails(
			http.StatusBadRequest,
			"Person-ID not correctly specified",
			"", nil))
		return
	}
	var absences []int
	if err := json.NewDecoder(r.Body).Decode(&absences); err != nil {
		apihelper.ResponseBadRequest(w, funcName, apimodel.Result{Result: "meetings id not valid"}, err)
		return
	}
	var absencePerson []dbModel.PersonRecurringAbsence
	for _, ab := range absences {
		absencePerson = append(absencePerson,
			dbModel.PersonRecurringAbsence{PersonID: uint(*idPerson), Weekday: ab})
	}

	db := database.DB.Begin()
	defer db.Commit()
	err = absence.DeleteRecurringAbsence(absencePerson, db)
	if err != nil {
		db.Rollback()
		switch err {
		case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusBadRequest)
		default:
			api_helper.InternalError(w, funcName, err.Error())
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
