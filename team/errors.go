package team

import "errors"

var (
	errEmptyID   = errors.New("Team ID must be present")
	errEmptyName = errors.New("Team name must be present")
	errNoTeams   = errors.New("No teams present")
	errNoTeamId  = errors.New("Team is not present")
)
