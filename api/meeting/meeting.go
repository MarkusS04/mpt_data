package meeting

import (
	"encoding/json"

	"mpt_data/api/apihelper"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/api/auth"
	"mpt_data/database/meeting"
	"mpt_data/helper"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

const packageName = "api.meeting"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc(apiModel.MeetingHref, auth.CheckAuthentication(getMeetings)).Methods(http.MethodGet)
	mux.HandleFunc(apiModel.MeetingHref, auth.CheckAuthentication(addMeeting)).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.MeetingHrefWithID, auth.CheckAuthentication(updatetMeeting)).Methods(http.MethodPut)
	mux.HandleFunc(apiModel.MeetingHrefWithID, auth.CheckAuthentication(deleteMeeting)).Methods(http.MethodDelete)
}

// @Summary		Get Meetings
// @Description	Get all Meetings in the specified time period
// @Tags			Meeting
// @Accept			json
// @Produce		json
// @Param			StartDate	query	string	true	"Start date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Param			EndDate		query	string	true	"End date/timestamp, Either English Date, or RFC3339"	Example("2023-01-21", "2023-01-21T00:00:00+00:00")
// @Security		ApiKeyAuth
// @Success		200 {array} dbModel.Meeting
// @Failure		400
// @Failure		401
// @Router			/meeting [GET]
func getMeetings(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".getMeetings"
	queryParams := r.URL.Query()

	startDate, err := helper.ParseTime(queryParams.Get("StartDate"))
	endDate, err2 := helper.ParseTime(queryParams.Get("EndDate"))
	if err != nil || err2 != nil {
		apihelper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "could not parse StartDate and/or EndDate"}, err)
		return
	}

	meetings, err := meeting.GetMeetings(generalmodel.Period{StartDate: startDate, EndDate: endDate})
	if err != nil {
		api_helper.InternalError(w, funcName, err.Error())
		return
	}
	api_helper.ResponseJSON(w, funcName, meetings)
}

// @Summary		Add Meetings
// @Description	Add Meetings
// @Tags			Meeting
// @Accept			json
// @Produce		json
// @Param			Meetings	body	[]dbModel.Meeting	true	"Meetings"
// @Security		ApiKeyAuth
// @Success		200
// @Success		201
// @Failure		400
// @Failure		401
// @Router			/meeting [POST]
func addMeeting(w http.ResponseWriter, r *http.Request) {
	var meetings []dbModel.Meeting
	if err := json.NewDecoder(r.Body).Decode(&meetings); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := meeting.AddMeetings(meetings)
	if err != nil {
		switch err {
		case errors.ErrNotAllMeetingsCreated:
			apihelper.ResponseJSON(w, packageName+"addMeeting", apiModel.Result{Result: "not all meetings created"})
		case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusBadRequest)
		default:
			api_helper.InternalError(w, packageName+"addMeeting", err.Error())
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary		Update Meetings
// @Description	Update the date of one meeting
// @Tags			Meeting
// @Accept			json
// @Produce		json
// @Param			id		path	int				true	"ID of meeting"
// @Param			Meeting	body	dbModel.Meeting	true	"Meeting"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400
// @Failure		401
// @Router			/meeting/{id} [PUT]
func updatetMeeting(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + "updateMeeting"
	var meetingIn dbModel.Meeting
	var err error
	if err = json.NewDecoder(r.Body).Decode(&meetingIn); err != nil {
		apihelper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "error in request body"}, err)

		return
	}

	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "id not valid"}, err)
		return
	}
	meetingIn.ID = uint(*id)

	err = meeting.UpdateMeeting(meetingIn)
	if err != nil {
		api_helper.InternalError(w, packageName+".updateMeeting", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary		Delete Meetings
// @Description	Delete one Meeting
// @Tags			Meeting
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"ID of meeting"
// @Security		ApiKeyAuth
// @Success		200
// @Failure		400
// @Failure		401
// @Router			/meeting/{id} [DELETE]
func deleteMeeting(w http.ResponseWriter, r *http.Request) {
	const funcName = packageName + ".deleteMeeting"
	var meetingIn dbModel.Meeting
	id, err := helper.ExtractIntFromURL(r, "id")
	if err != nil || *id <= 0 {
		apihelper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "id not valid"}, err)
		return
	}
	meetingIn.ID = uint(*id)

	if err := meeting.DeleteMeeting(meetingIn); err == errors.ErrMeetingNotDeleted {
		apihelper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "meetings not deleted"}, err)
	} else if err != nil {
		api_helper.InternalError(w, packageName+".deleteMeeting", err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
