package team

import (
	"encoding/json"
	"net/http"

	"github.com/A9u/function_junction/api"
  "github.com/gorilla/mux"
  "github.com/mongodb/mongo-go-driver/bson/primitive"
)

func Create(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		eventID, err := primitive.ObjectIDFromHex(vars["event_id"])
		var c createRequest
		err = json.NewDecoder(req.Body).Decode(&c)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

    response, err := service.create(req.Context(), c, eventID)
		if isBadRequest(err) {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusCreated, response)
	})
}

func List(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		eventID, err := primitive.ObjectIDFromHex(vars["event_id"])
		resp, err := service.list(req.Context(), eventID)
		if err == errNoTeams {
			api.Error(rw, http.StatusNotFound, api.Response{Message: err.Error()})
			return
		}
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, resp)
	})
}

func isBadRequest(err error) bool {
	return err == errEmptyName || err == errEmptyID
}
