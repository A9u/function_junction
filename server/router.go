package server

import (
	"fmt"
	"net/http"
	"github.com/A9u/function_junction/api"
	"github.com/A9u/function_junction/config"
	"github.com/A9u/function_junction/team"
	"github.com/gorilla/mux"
	"github.com/A9u/function_junction/event"
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

	router.HandleFunc("/events/{event_id}/teams", team.Create(dep.TeamService)).Methods(http.MethodPost).Headers(versionHeader, v1)
	router.HandleFunc("/events/{event_id}/teams", team.List(dep.TeamService)).Methods(http.MethodGet).Headers(versionHeader, v1)

	sh := http.StripPrefix("/docs/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/docs/").Handler(sh)
	return
}

func pingHandler(rw http.ResponseWriter, req *http.Request) {
	api.Success(rw, http.StatusOK, api.Response{Message: "pong"})
}
