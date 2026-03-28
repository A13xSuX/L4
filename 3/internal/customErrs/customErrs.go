package customErrs

import "errors"

var (
	InvalidIDErr              = errors.New("invalid id")
	InvalidUserIDErr          = errors.New("invalid UserID")
	EventPastTimeErr          = errors.New("eventTime must be in the future")
	TitleEmptyErr             = errors.New("title empty")
	RemindAtPastErr           = errors.New("remindAt must be in the future")
	RemindAtAfterEventTimeErr = errors.New("remindAt must be before eventTime")
	PriorityErr               = errors.New("priority not in list: high, medium, low")
)
