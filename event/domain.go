package event

import (
	"github.com/A9u/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"time"
)

type updateRequest struct {
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	StartDateTime     time.Time     `json:"start_date_time"`
	EndDateTime       time.Time     `json:"end_date_time"`
	IsShowcasable     bool          `json:"is_showcasable"`
	IsIndividualEvent bool          `json:"is_individual_participation"`
	MaxSize           int           `json:"max_size"`
	MinSize           int           `json:"min_size"`
	IsPublished       bool          `json:"is_published"`
	Venue             string        `json:"venue"`
	UpdatedAt         time.Time     `json:"updated_at"`
	RegisterBefore    time.Time     `json:"register_before"`
	Set  db.Event `json:"$set"`
}

type createRequest struct {
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	StartDateTime     time.Time     `json:"start_date_time"`
	EndDateTime       time.Time     `json:"end_date_time"`
	IsShowcasable     bool          `json:"is_showcasable"`
	IsIndividualEvent bool          `json:"is_individual_participation"`
	CreatedBy         primitive.ObjectID  `json:"created_by" bson:"created_by"`
	MaxSize           int           `json:"max_size"`
	MinSize           int           `json:"min_size"`
	IsPublished       bool          `json:"is_published"`
	Venue             string        `json:"venue"`
	CreatedAt         time.Time     `json:"created_at"`
	RegisterBefore    time.Time     `json:"register_before"`
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
