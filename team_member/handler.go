package team_member

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joshsoftware/function_junction/api"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"net/http"
)

func Create(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var c createRequest

		queryParams := mux.Vars(req)
		teamID, err := primitive.ObjectIDFromHex(queryParams["team_id"])
		eventID, err1 := primitive.ObjectIDFromHex(queryParams["event_id"])
		fmt.Println("recieved params teamid", teamID)
		err = json.NewDecoder(req.Body).Decode(&c)
		if err != nil || err1 != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		resp, err := service.create(req.Context(), c, teamID, eventID)
		if isBadRequest(err) {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusCreated, resp)
	})
}

func List(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Println("team", vars["team_id"])
		id, err := primitive.ObjectIDFromHex(vars["team_id"])
		eventID, err := primitive.ObjectIDFromHex(vars["event_id"])
		fmt.Println("id", id)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		resp, err := service.list(req.Context(), id, eventID)
		if err == errNoTeamMember {
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

func FindListOfInviters(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Println("event", vars["event_id"])
		eventID, err := primitive.ObjectIDFromHex(vars["event_id"])
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		resp, err := service.findListOfInviters(req.Context(), eventID)
		if err == errNoTeamMember {
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

func FindByID(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Println(vars)
		id, err := primitive.ObjectIDFromHex(vars["team_member_id"])
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
		}

		resp, err := service.findByID(req.Context(), id)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
		}
		api.Success(rw, http.StatusOK, resp)
	})
}

func DeleteByID(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := primitive.ObjectIDFromHex(vars["team_member_id"])
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
		}

		err = service.deleteByID(req.Context(), id)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, api.Response{Message: "Deleted Successfully"})
	})
}

func Update(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		teamMemberID, err := primitive.ObjectIDFromHex(vars["team_member_id"])
		teamID, err := primitive.ObjectIDFromHex(vars["team_id"])
		eventId, err := primitive.ObjectIDFromHex(vars["event_id"])

		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
		}

		var c updateRequest
		err = json.NewDecoder(req.Body).Decode(&c)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
			return
		}

		resp, err := service.update(req.Context(), c, teamMemberID, teamID, eventId)
		if isBadRequest(err) {
			api.Error(rw, http.StatusBadRequest, api.Response{Message: err.Error()})
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
	return err == errEmptyID || err == errEmptyEmail
}

func CancelRsvp(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		eventId, err := primitive.ObjectIDFromHex(vars["event_id"])
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		err = service.cancelRsvp(req.Context(), eventId)
		if err != nil {
			api.Error(rw, http.StatusInternalServerError, api.Response{Message: err.Error()})
			return
		}

		api.Success(rw, http.StatusOK, api.Response{Message: "Cancelled Rsvp successfully!"})
	})
}
