package team

import (
	"time"

	"github.com/A9u/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type createRequest struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	EventID     primitive.ObjectID `json:"event_id"`
	CreatorID   primitive.ObjectID `json:"creator_id"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	ShowcaseUrl string             `json:"showcase_url"`
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
