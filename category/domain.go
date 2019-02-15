package category

import "github.com/joshsoftware/golang-boilerplate/db"

type updateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Set  db.Category `json:"$set"`
}

type createRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type findByIDResponse struct {
	Category db.Category `json:"category"`
}

type listResponse struct {
	Categories []*db.Category `json:"categories"`
}

func (cr createRequest) Validate() (err error) {
	if cr.Name == "" {
		return errEmptyName
	}
	return
}

func (ur updateRequest) Validate() (err error) {
	if ur.Name == "" {
		return errEmptyName
	}
	return
}
