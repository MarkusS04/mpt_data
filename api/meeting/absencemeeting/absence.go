package absencemeeting

import (
	"encoding/json"
	"mpt_data/api/apihelper"
	"mpt_data/api/middleware"
	"mpt_data/database/absence"
	"mpt_data/helper"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

const packageName = "api.meeting.absence"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc(apiModel.MeetingAbsence, middleware.CheckAuthentication(getAbsence)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.MeetingAbsence, middleware.CheckAuthentication(addAbsence)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.MeetingAbsence, middleware.CheckAuthentication(deleteAbsence)).Methods(http.MethodDelete)
}

// @Summary		Get Absence
// @Description	Get absent people for a meeting
// @Tags			Meeting,Absence
// @Accept			json
// @Produce		json
// @Param			MeetingId	path	int	true	"ID of meeting"
// @Security		ApiKeyAuth
// @Success		200 {object} dbModel.Person
// @Failure		400
// @Failure		401
// @Router			/meeting/{MeetingId}/absence [GET]
func getAbsence(w http.ResponseWriter, r *http.Request) {
	var id *int
	var err error
	if id, err = helper.ExtractIntFromURL(r, "id"); err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, apiModel.Result{Result: "id not valid"}, err)
		return
	}

	data, err := absence.GetAbsenceMeeting(uint(*id))
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusOK)
		case gorm.ErrEmptySlice, gorm.ErrInvalidData:
			w.WriteHeader(http.StatusBadRequest)
		default:
			apihelper.InternalError(w, err)
		}
	} else {
		apihelper.ResponseJSON(w, data)
	}
}

// @Summary		Add Absence
// @Description	Add absent people for a meeting
// @Tags			Meeting,Absence
// @Accept			json
// @Produce		json
// @Param			MeetingId	path	int		true	"ID of meeting"
// @Param			Absence		body	[]uint	true	"ID of people who are absent"
// @Security		ApiKeyAuth
// @Success		200
// @Success		201
// @Failure		400
// @Failure		401
// @Router			/meeting/{MeetingId}/absence [POST]
func addAbsence(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".addAbsence"
	var absences []uint
	var absencePerson []dbModel.PersonAbsence
	if err := json.NewDecoder(r.Body).Decode(&absences); err != nil {
		apihelper.ResponseBadRequest(w, apiModel.Result{Result: "people id not set"}, err)
		return
	}

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, apiModel.Result{Result: "meeting id not valid"}, err)
		return
	}
	for _, ab := range absences {
		absencePerson = append(absencePerson,
			dbModel.PersonAbsence{MeetingID: uint(*id), PersonID: ab})
	}

	err = absence.AddAbsence(absencePerson)
	if err != nil {
		switch err {
		case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusBadRequest)
		default:
			apihelper.InternalError(w, err)
		}
	} else {
		apihelper.ResponseJSON(w, absences, http.StatusCreated)
	}
}

// @Summary		Delete Absence
// @Description	Delete absent people for a meeting
// @Tags			Meeting,Absence
// @Accept			json
// @Produce		json
// @Param			MeetingId	path	int		true	"ID of meeting"
// @Param			Absence		body	[]uint	true	"ID of people who are no longer absent"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400
// @Failure		401
// @Router			/meeting/{MeetingId}/absence [DELETE]
func deleteAbsence(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deleteAbsence"
	var absences []uint
	var absencePerson []dbModel.PersonAbsence
	if err := json.NewDecoder(r.Body).Decode(&absences); err != nil {
		apihelper.ResponseBadRequest(w, apiModel.Result{Result: "people id not set"}, err)
		return
	}

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, apiModel.Result{Result: "id not valid"}, err)
		return
	}
	for _, ab := range absences {
		absencePerson = append(absencePerson,
			dbModel.PersonAbsence{MeetingID: uint(*id), PersonID: ab})
	}

	err = absence.DeleteAbsence(absencePerson)
	if err != nil {
		switch err {
		case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusBadRequest)
		default:
			apihelper.InternalError(w, err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
