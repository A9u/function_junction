package event

import (
	"github.com/A9u/function_junction/db"
	"time"
)

type createRequest struct {
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
	RegisterBefore    time.Time     `json:"register_before"`
}

type updateRequest struct {
	createRequest
	Set  db.Event `json:"$set"`
}

type listResponse struct {
	Events []*db.EventInfo `json:"events"`
}

type eventResponse struct {
	Event *db.EventInfo `json:"event"`
}

func (cr createRequest) EventValidate() (err error) {
	if cr.Title == "" {
		return errEmptyTitle
	}

	if cr.EndDateTime.IsZero() || cr.StartDateTime.IsZero() {
		return errEmptyDate
	}

	if cr.IsIndividualEvent == false && (cr.MaxSize == 0 || cr.MinSize == 0){
		return errEmptyTeamSize
	}

	if cr.IsIndividualEvent == false && (cr.MinSize > cr.MaxSize){
		return errInvalidTeamSize
	}
	return
}
