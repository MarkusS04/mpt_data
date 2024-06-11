// Package apihelper provides general functions to use in API
package apihelper

import (
	"encoding/json"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	generalmodel "mpt_data/models/general"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

/*
InternalError logs error and sends http.StatusInternalServerError

	err: text, that will be logged
	both not send to client
*/
func InternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	zap.L().Error(generalmodel.InternalError, zap.Error(err))
}

/*
ResponseJSON sends specified data to client as json

	tries to translates data into json
	logs if any error occurs
	if statusCode is set, uses statusCode, otherwise http.StatusOK
*/
func ResponseJSON(w http.ResponseWriter, data interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if json, err := json.Marshal(data); err != nil {
		InternalError(w, err)
	} else {
		w.WriteHeader(code)
		if _, err := w.Write(json); err != nil {
			InternalError(w, err)
		}
	}
}

/*
ResponseError sends specified data to client as problem+json

	tries to translates data into json
	logs if any error occurs
	uses statusCode of data
*/

func ResponseError(w http.ResponseWriter, data apiModel.ProblemDetails) {
	code := data.Status

	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	if json, err := json.Marshal(data); err != nil {
		InternalError(w, err)
	} else {
		w.WriteHeader(code)
		if _, err := w.Write(json); err != nil {
			InternalError(w, err)
		}
	}
}

/*
ResponseBadRequest sends data to client with StatusBadRequest

	logs the error with data
*/
func ResponseBadRequest(w http.ResponseWriter, data apiModel.Result, err error) {
	zap.L().Error(generalmodel.StatusBadRequest, zap.Any("Result", data.Result), zap.Error(err))
	ResponseJSON(w, data, http.StatusBadRequest)
}

func ExtractIntFromURL(r *http.Request, fieldName string) (int64, error) {
	vars := mux.Vars(r)
	field, exists := vars[fieldName]

	if !exists {
		// Field not in URL
		return 0, errors.ErrPathParamMissing
	}
	fieldInt, err := strconv.ParseInt(field, 10, 0)
	if err != nil {
		// error parsing field
		return 0, errors.ErrPathWrongType
	}

	return fieldInt, nil
}
