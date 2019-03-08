package team

import (
	"time"

	"github.com/A9u/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type createRequest struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	EventID     primitive.ObjectID `json:"eventId"`
	CreatorID   primitive.ObjectID `json:"creatorId"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	ShowcaseUrl string             `json:"showcaseUrl"`
	Description string             `json:"description"`
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

type createResponse struct {
  Team *db.Team `json:"team"`
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
