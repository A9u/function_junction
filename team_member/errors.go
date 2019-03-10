package team_member

import "errors"

var (
	errEmptyEmail     = errors.New("TeamMember email must be present")
	errEmptyEmails    = errors.New("TeamMember emails must be present")
	errEmptyID      = errors.New("TeamMember ID must be present")
	errEmptyName    = errors.New("TeamMember name must be present")
	errEmptyInviterMail    = errors.New("Inviter Email ID must be present")
	errNoTeamMember = errors.New("No TeamMember present")
	errNoTeamMemberId = errors.New("TeamMember is not present")
	errInviterDoesNotExist = errors.New("Inviter Does not Exist in Db")
	errEventDoesNotExist = errors.New("Event Does not Exist in Db")
	errTeamDoesNotExist = errors.New("Team Does not Exist in Db")

)
