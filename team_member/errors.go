package team_member

import "errors"

var (
	errEmptyID        = errors.New("TeamMember ID must be present")
	errEmptyName      = errors.New("TeamMember name must be present")
	errNoTeamMember   = errors.New("No TeamMember present")
	errNoTeamMemberId = errors.New("TeamMember is not present")
	errEmptyEmail     = errors.New("TeamMember email must be present")
	errEmptyEmails    = errors.New("TeamMember emails must be present")
)
