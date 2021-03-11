package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/function_junction/api"
	"github.com/joshsoftware/function_junction/config"
	"github.com/joshsoftware/function_junction/event"
	"github.com/joshsoftware/function_junction/team"
	"github.com/joshsoftware/function_junction/team_member"
)

const (
	versionHeader = "Accept"
)

func initRouter(dep dependencies) (router *mux.Router) {
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	router = mux.NewRouter()
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	router.HandleFunc("/events", event.Create(dep.EventService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/events", event.List(dep.EventService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}", event.FindByID(dep.EventService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}", event.DeleteByID(dep.EventService)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}", event.Update(dep.EventService)).Methods(http.MethodPut).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/rsvp", team_member.Create(dep.TeamMemberService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/{email}/rsvp/no", team_member.Reject(dep.TeamMemberService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/{email}/rsvp/yes", team_member.Accept(dep.TeamMemberService)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.HandleFunc("/events/{event_id}/teams", team.Create(dep.TeamService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams", team.List(dep.TeamService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams/{team_id}", team.DeleteByID(dep.TeamService)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams/{team_id}", team.Update(dep.TeamService)).Methods(http.MethodPut).Headers(versionHeader, v1)

	// TeamMember
	router.HandleFunc("/events/{event_id}/teams/{team_id}/team_members", team_member.Create(dep.TeamMemberService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams/{team_id}/team_members", team_member.List(dep.TeamMemberService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/attendees", team_member.List(dep.TeamMemberService)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.HandleFunc("/events/{event_id}/teams/{team_id}/team_members/{team_member_id}", team_member.FindByID(dep.TeamMemberService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams/{team_id}/team_members/{team_member_id}", team_member.DeleteByID(dep.TeamMemberService)).Methods(http.MethodDelete).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams/{team_id}/team_members/{team_member_id}", team_member.Update(dep.TeamMemberService)).Methods(http.MethodPut).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/invited_by/", team_member.FindListOfInviters(dep.TeamMemberService)).Methods(http.MethodGet).Headers(versionHeader, v1)
	sh := http.StripPrefix("/docs/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/docs/").Handler(sh)
	return
}

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	api.Success(rw, http.StatusOK, api.Response{Message: "pong"})
}
