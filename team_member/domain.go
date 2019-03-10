package team_member

import (
	// "regexp"
	"github.com/A9u/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type updateRequest struct {
	Status    string             `json:"status"`
}

type createRequest struct {
	Emails []string `json:emails`
}

type InviterslistResponse struct {
	InvitersInfo []*db.InviterInfo `json:"Inviters"`
}

type findByIDResponse struct {
	TeamMember db.TeamMember `json:"team_member"`
}

type listResponse struct {
	TeamMembers []*db.TeamMemberInfo `json:"team_members"`
}

type createResponse struct {
	FailedEmails []string `json:"failed_emails"`
}

type updateResponse struct {
	TeamMember db.TeamMember `json:"team_member"`
}

func (cr createRequest) Validate() (err error) {
	if len(cr.Emails) == 0 {
		return errEmptyEmails
	}
	return
}

func (ur updateRequest) Validate() (err error) {
	if ur.Status == "" {
		return errEmptyStatus
	}

	// re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@joshsoftware.com$")

	return
}
