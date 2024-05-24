package apihelper

import (
	"encoding/json"
	"fmt"
	"mpt_data/database/logging"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const packageName = "apihelper"

/*
InternalError logs error and sends http.StatusInternalServerError

	funcName: Name for logging
	err: text, that will be logged
	both not send to client
*/
func InternalError(w http.ResponseWriter, funcName string, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	logging.LogError(funcName, err)
}

/*
ResponseJSON sends specified data to client as json

	tries to translates data into json
	logs if any error occurs
	if statusCode is set, uses statusCode, otherwise http.StatusOK
*/
func ResponseJSON(w http.ResponseWriter, funcName string, data interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if json, err := json.Marshal(data); err != nil {
		InternalError(w, funcName, err.Error())
	} else {
		w.WriteHeader(code)
		if _, err := w.Write(json); err != nil {
			logging.LogError(packageName+"ResponseJSON", err.Error())
		}
	}
}

/*
ResponseError sends specified data to client as problem+json

	tries to translates data into json
	logs if any error occurs
	uses statusCode of data
*/

func ResponseError(w http.ResponseWriter, funcName string, data apiModel.ProblemDetails) {
	code := data.Status

	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	if json, err := json.Marshal(data); err != nil {
		InternalError(w, funcName, err.Error())
	} else {
		w.WriteHeader(code)
		if _, err := w.Write(json); err != nil {
			logging.LogError(packageName+"ResponseError", err.Error())
		}
	}
}

/*
ResponseBadRequest sends data to client with StatusBadRequest

	logs the error with data
*/
func ResponseBadRequest(w http.ResponseWriter, funcName string, data apiModel.Result, err error) {
	logging.LogError(funcName, fmt.Sprintf("%v: %v", data.Result, err))
	ResponseJSON(w, funcName, data, http.StatusBadRequest)
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
