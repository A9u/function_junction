package team_member

import "github.com/A9u/function_junction/db"
import "github.com/mongodb/mongo-go-driver/bson/primitive"

type updateRequest struct {
	Name      string             `json:"name"`
	Status    string             `json:"status"`
	InviterID primitive.ObjectID `json:"inviter_id"`
	TeamID    primitive.ObjectID `json:"team_id"`
	EventID   primitive.ObjectID `json:"event_id"`
}

type createRequest struct {
	Emails []string
	TeamID primitive.ObjectID `json:"team_id"`
}

type findByIDResponse struct {
	TeamMember db.TeamMember `json:"team_member"`
}

type listResponse struct {
	TeamMembers []*db.TeamMember `json:"team_members"`
}

func (cr createRequest) Validate() (err error) {
	if len(cr.Emails) == 0 {
		return errEmptyEmails
	}
	return
}

func (ur updateRequest) Validate() (err error) {
	if ur.Name == "" {
		return errEmptyName
	}
	return
}
