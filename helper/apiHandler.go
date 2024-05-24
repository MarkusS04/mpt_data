package helper

import (
	"mpt_data/helper/errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ExtractIntFromURL(r *http.Request, fieldName string) (*int, error) {
	vars := mux.Vars(r)
	field, exists := vars[fieldName]

	if !exists {
		// Field not in URL
		return nil, errors.ErrPathParamMissing
	}
	fieldInt, err := strconv.ParseInt(field, 10, 0)
	if err != nil {
		// error parsing field
		return nil, errors.ErrPathWrongType
	}

	value := int(fieldInt)
	return &value, nil
}
