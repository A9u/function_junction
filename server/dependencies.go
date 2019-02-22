package server

import (
	"github.com/A9u/function_junction/app"
	"github.com/A9u/function_junction/category"
	"github.com/A9u/function_junction/event"
	"github.com/A9u/function_junction/db"
	"github.com/A9u/function_junction/team"
)

type dependencies struct {
	CategoryService category.Service
	TeamService     team.Service
	EventService event.Service
}

func initDependencies() (dependencies, error) {
	appDB := app.GetDB()
	logger := app.GetLogger()
	dbStore := db.NewStorer(appDB)

	categoryService := category.NewService(dbStore, logger, app.GetCollection("catagories"))
	teamService := team.NewService(dbStore, logger, app.GetCollection("teams"))
	eventService := event.NewService(dbStore, logger, app.GetCollection("events"))o

  return dependencies{
		CategoryService: categoryService,
		EventService: eventService,
		TeamService:     teamService,
	}, nil
}
