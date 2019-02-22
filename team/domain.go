package team

import "github.com/A9u/function_junction/db"
import "github.com/mongodb/mongo-go-driver/bson/primitive"

type createRequest struct {
	Name    string             `json:"name"`
	Type    string             `json:"type"`
	EventID primitive.ObjectID `json:"event_id"`
}

type listResponse struct {
	Teams []*db.Team `json:"teams"`
}

func (cr createRequest) Validate() (err error) {
	if cr.Name == "" {
		return errEmptyName
	}
	return
}
