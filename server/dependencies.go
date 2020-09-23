package server

import (
	"github.com/joshsoftware/function_junction/app"
	"github.com/joshsoftware/function_junction/category"
	"github.com/joshsoftware/function_junction/db"
	"github.com/joshsoftware/function_junction/event"
	"github.com/joshsoftware/function_junction/mailer"
	"github.com/joshsoftware/function_junction/team"
	"github.com/joshsoftware/function_junction/team_member"
)

type dependencies struct {
	CategoryService   category.Service
	TeamService       team.Service
	EventService      event.Service
	TeamMemberService team_member.Service
	MailerService     mailer.MailerService
}

func initDependencies() (dependencies, error) {
	appDB := app.GetDB()
	logger := app.GetLogger()
	dbStore := db.NewStorer(appDB)

	mailerService := mailer.NewMailer()
	categoryService := category.NewService(dbStore, logger, app.GetCollection("catagories"))
	teamService := team.NewService(dbStore, logger, app.GetCollection("teams"))
	eventService := event.NewService(dbStore, logger, app.GetCollection("events"), mailerService)
	teamMemberService := team_member.NewService(dbStore, logger, app.GetCollection("team_members"),
		app.GetCollection("teams"), app.GetCollection("users"), app.GetCollection("events"), mailerService)

	return dependencies{
		CategoryService:   categoryService,
		EventService:      eventService,
		TeamService:       teamService,
		TeamMemberService: teamMemberService,
		MailerService:     mailerService,
	}, nil
}
