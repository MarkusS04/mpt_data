package apimodel

type Result struct {
	Result string
}

type ProblemDetails struct {
	Type     string      `json:"type,omitempty"`
	Title    string      `json:"title,omitempty"`
	Status   int         `json:"status,omitempty"`
	Detail   string      `json:"detail,omitempty"`
	Instance string      `json:"instance,omitempty"`
	Errors   interface{} `json:"information,omitempty"`
}

func GetProblemDetails(typeVal, title string, status int, detail, instance string, errors interface{}) *ProblemDetails {
	return &ProblemDetails{
		Type:     typeVal,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
		Errors:   errors,
	}
}

func GetInavalidRequestProblemDetails(status int, detail, instance string, errors interface{}) *ProblemDetails {
	return &ProblemDetails{
		Type:     "error:invalid-request",
		Title:    "Your request is not valid",
		Status:   status,
		Detail:   detail,
		Instance: instance,
		Errors:   errors,
	}
}
