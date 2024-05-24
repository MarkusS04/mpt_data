package apimodel

const (
	base = "/api/v1"
)

const (
	LoginHref        = base + "/login"
	UserHref         = base + "/user"
	UserChangePWHref = UserHref + "/password"
)
const (
	MeetingHref       = base + "/meeting"
	MeetingHrefWithID = MeetingHref + "/{id}"
)

const (
	PlanHref             = base + "/plan"
	PlanHrefPDF          = PlanHref + "/pdf"
	PlanHrefWithID       = PlanHref + "/{id}"
	PlanHrefWithIDPeople = PlanHrefWithID + "/people"
)

const (
	MeetingAbsence              = MeetingHrefWithID + "/absence"
	PersonAbsence               = PersonHrefWithID + "/absence"
	PersonAbsenceRecuring       = PersonHrefWithID + "/absencerecurring"
	PersonAbsenceRecuringWithID = PersonAbsenceRecuring + "/{absenceId}"
)

const (
	TaskHref             = base + "/task"
	TaskHrefWithID       = TaskHref + "/{id}"
	TaskDetailHref       = TaskHrefWithID + "/detail"
	TaskDetailHrefWithID = TaskDetailHref + "/{detailId}"
)

const (
	PersonHref       = base + "/person"
	PersonHrefWithID = PersonHref + "/{id}"
)

const (
	PersonHrefTask = PersonHrefWithID + "/task"
)
