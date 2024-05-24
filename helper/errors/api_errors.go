package errors

import "errors"

var (
	ErrPathParamMissing = errors.New("parameter not in path")
	ErrPathWrongType    = errors.New("path parameter of wrong type")
)
