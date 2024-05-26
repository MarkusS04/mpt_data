package apimodel

const (
	base = "/api/v1"
)

// User Routes for API
const (
	LoginHref        = base + "/login"
	UserHref         = base + "/user"
	UserChangePWHref = UserHref + "/password"
)

// Meeting Routes for API
const (
	MeetingHref       = base + "/meeting"
	MeetingHrefWithID = MeetingHref + "/{id}"
	MeetingTagHref    = MeetingHrefWithID + "/tag"
)

// Plan Routes for API
const (
	PlanHref             = base + "/plan"
	PlanHrefPDF          = PlanHref + "/pdf"
	PlanHrefWithID       = PlanHref + "/{id}"
	PlanHrefWithIDPeople = PlanHrefWithID + "/people"
)

// Absence Routes for API
const (
	MeetingAbsence              = MeetingHrefWithID + "/absence"
	PersonAbsence               = PersonHrefWithID + "/absence"
	PersonAbsenceRecuring       = PersonHrefWithID + "/absencerecurring"
	PersonAbsenceRecuringWithID = PersonAbsenceRecuring + "/{absenceId}"
)

// Task Routes for API
const (
	TaskHref             = base + "/task"
	TaskHrefWithID       = TaskHref + "/{id}"
	TaskDetailHref       = TaskHrefWithID + "/detail"
	TaskDetailHrefWithID = TaskDetailHref + "/{detailId}"
)

// Person Routes for API
const (
	PersonHref       = base + "/person"
	PersonHrefWithID = PersonHref + "/{id}"
	PersonHrefTask   = PersonHrefWithID + "/task"
)
