package errors

import (
	"errors"
)

var (
	ErrUserdataNotComplete = errors.New("username or hash not set")
	ErrUserNotComplete     = errors.New("username or password not set")

	ErrInvalidAuth = errors.New("incorect password or username")
	ErrFailedAuth  = errors.New("error in authentication")
)
