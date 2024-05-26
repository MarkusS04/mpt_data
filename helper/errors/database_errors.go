package errors

import (
	"errors"
)

var ErrUserAlreadyExists = errors.New("user already exists")

var (
	ErrNotAllMeetingsCreated = errors.New("not all given meetings written to DB")
	ErrMeetingNotDeleted     = errors.New("meeting not deleted")
	ErrMeetingTagAlreadySet  = errors.New("tag for meeting is already set")
)

var (
	ErrTaskForPersonNotAllowed = errors.New("person is not allowed for task")
)

var (
	ErrIDNotSet        = errors.New("the id was not set")
	ErrForeignIDNotSet = errors.New("the foreign id was not set")
)
