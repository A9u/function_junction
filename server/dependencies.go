package server

import (
	"github.com/A9u/function_junction/app"
	"github.com/A9u/function_junction/category"
	"github.com/A9u/function_junction/team_member"
	"github.com/A9u/function_junction/db"
)

type dependencies struct {
	CategoryService category.Service
	TeamMemberService team_member.Service
}

func initDependencies() (dependencies, error) {
	appDB := app.GetDB()
	logger := app.GetLogger()
	dbStore := db.NewStorer(appDB)

	categoryService := category.NewService(dbStore, logger, app.GetCollection("catagories"))
	teamMemberService := team_member.NewService(dbStore, logger, app.GetCollection("team_members"))

	return dependencies{
		CategoryService: categoryService,
		TeamMemberService: teamMemberService,
	}, nil
}
