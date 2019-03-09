package event

import "errors"

var (
	errEmptyID 				= errors.New("Event ID must be present")
	errEmptyTitle 			= errors.New("Event name must be present")
	errNoEvents 			= errors.New("No events present")
	errNoEventId 			= errors.New("Event is not present")
	errEmptyDate 		= errors.New("Event Start and End Date and Time must be present")
	errEmptyIndividualEvent = errors.New("Please Select if event is Individual or not")
	errInvalidStartDate 	= errors.New("Event Start Date and Time must be present")
	errInvalidEndDate 		= errors.New("Event End Date and Time must be present")
	errInvalidLastDate 		= errors.New("Event End Date and Time must be present")
	errEmptyTeamSize				= errors.New("Minimum and maximum team size are compulsory for team events")
	errInvalidTeamSize				= errors.New("Minimum Team size should be greater than Minimum tesm size")
)
