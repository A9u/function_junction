package team

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/A9u/function_junction/db"
	"go.uber.org/zap"
)

type Service interface {
	list(ctx context.Context) (response listResponse, err error)
	create(ctx context.Context, req createRequest) (err error)
}

type teamService struct {
	store      db.Storer
	logger     *zap.SugaredLogger
	collection *mongo.Collection
}

func (ts *teamService) list(ctx context.Context) (response listResponse, err error) {
	teams, err := ts.store.ListTeams(ctx, ts.collection)
	fmt.Println("========")
	fmt.Println(teams)
	if err == db.ErrTeamNotExist {
		ts.logger.Error("No team present", "err", err.Error())
		return response, errNoTeams
	}
	if err != nil {
		ts.logger.Error("Error listing teams", "err", err.Error())
		return
	}

	response.Teams = teams
	return
}

func (ts *teamService) create(ctx context.Context, c createRequest) (err error) {
	err = c.Validate()
	if err != nil {
		ts.logger.Errorw("Invalid request for team create", "msg", err.Error(), "team", c)
		return
	}

	err = ts.store.CreateTeam(ctx, ts.collection, &db.Team{Name: c.Name})
	if err != nil {
		ts.logger.Error("Error creating team", "err", err.Error())
		return
	}
	return
}

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection) Service {
	return &teamService{
		store:      s,
		logger:     l,
		collection: c,
	}
}
