package event

import (
	"github.com/A9u/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"time"
	)

type updateRequest struct {
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	StartDateTime     time.Time     `json:"startDateTime"`
	EndDateTime       time.Time     `json:"endDateTime"`
	IsShowcasable     bool          `json:"isShowcasable"`
	IsIndividualEvent bool          `json:"isIndividualParticipation"`
	MaxSize           int           `json:"maxSize"`
	MinSize           int           `json:"minSize"`
	IsPublished       bool          `json:"isPublished"`
	Venue             string        `json:"venue"`
	UpdatedAt         time.Time     `db:"updated_at"`
	Set  db.Event `json:"$set"`
}

type createRequest struct {
	Id                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	StartDateTime     time.Time     `json:"startDateTime"`
	EndDateTime       time.Time     `json:"endDateTime"`
	IsShowcasable     bool          `json:"isShowcasable"`
	IsIndividualEvent bool          `json:"isIndividualParticipation"`
	CreatedBy         primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	MaxSize           int           `json:"maxSize"`
	MinSize           int           `json:"minSize"`
	IsPublished       bool          `json:"isPublished"`
	Venue             string        `json:"venue"`
	CreatedAt         time.Time     `db:"createdAt"`
	UpdatedAt         time.Time     `db:"updatedAt"`
	RegisterBefore    time.Time     `db:"registerBefore"`
}

type findByIDResponse struct {
	Event db.Event `json:"event"`
}

type listResponse struct {
	Events []*db.Event `json:"events"`
}

func (cr createRequest) Validate() (err error) {
	if cr.Title == "" {
		return errEmptyTitle
	}
	return
}

type createResponse struct {
	Event *db.Event `json:"event"`
}

type updateResponse struct {
	Event *db.Event `json:"event"`
}