package team

import (
	"github.com/A9u/function_junction/db"
)

type createRequest struct {
	Name        string `json:"name"`
	ShowcaseUrl string `json:"showcase_url"`
	Description string `json:"description"`
}

type listResponse struct {
	Teams []*db.TeamInfo `json:"teams"`
}

func (cr createRequest) Validate() (err error) {
	if cr.Name == "" {
		return errEmptyName
	}
	return
}

type createResponse struct {
	Team *db.TeamInfo `json:"team"`
}

/*
func (cr createRequest) ValidateEvent(ctx, eventID) (err error) {
  _, err := db.store.FindEventByID(ctx, eventID, app.GetCollection("events")
  if err != "" {
    return errInvalidEventID
  }
  return
}
*/
