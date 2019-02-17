package event

import "errors"

var (
	errEmptyID      = errors.New("Event ID must be present")
	errEmptyTitle    = errors.New("Event name must be present")
	errNoEvents = errors.New("No events present")
	errNoEventId = errors.New("Event is not present")
)
