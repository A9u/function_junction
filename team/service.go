package team

import (
	"context"

	"github.com/joshsoftware/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
)

type Service interface {
	list(ctx context.Context, eventID primitive.ObjectID) (response listResponse, err error)
	create(ctx context.Context, req createRequest, eventID primitive.ObjectID) (response createResponse, err error)
	deleteByID(ctx context.Context, teamID primitive.ObjectID) (err error)
	update(ctx context.Context, req createRequest, teamID primitive.ObjectID) (response createResponse, err error)
}

type teamService struct {
	store      db.Storer
	logger     *zap.SugaredLogger
	collection *mongo.Collection
}

func (ts *teamService) list(ctx context.Context, eventID primitive.ObjectID) (response listResponse, err error) {
	teams, err := ts.store.ListTeams(ctx, ts.collection, eventID)
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

func (ts *teamService) create(ctx context.Context, c createRequest, eventID primitive.ObjectID) (response createResponse, err error) {
	err = c.Validate()
	if err != nil {
		ts.logger.Errorw("Invalid request for team create", "msg", err.Error(), "team", c)
		return
	}

	createdTeam, err := ts.store.CreateTeam(ctx, ts.collection, &db.Team{
		Name:        c.Name,
		Description: c.Description,
		ShowcaseUrl: c.ShowcaseUrl,
		EventID:     eventID,
		CreatorID:   ctx.Value("currentUser").(db.User).ID,
	})
	if err != nil {
		ts.logger.Error("Error creating team", "err", err.Error())
		return
	}
	response.Team = createdTeam
	return
}

func (ts *teamService) deleteByID(ctx context.Context, id primitive.ObjectID) (err error) {
	err = ts.store.DeleteTeamByID(ctx, id, ts.collection)
	if err != nil {
		ts.logger.Error("Error deleting Team - ", "err", err.Error(), "team_id", id)
		return
	}

	err = ts.store.DeleteAllTeamMembers(ctx, id)

	if err != nil {
		ts.logger.Error("Error deleting Team Members- ", "err", err.Error(), "team_id", id)
		return
	}

	return
}

func (ts *teamService) update(ctx context.Context, req createRequest, id primitive.ObjectID) (response createResponse, err error) {
	currentUserID := ctx.Value("currentUser").(db.User).ID
	_, err = ts.store.FindTeamMemberByInviteeIDTeamID(ctx, currentUserID, id)

	if err != nil {
		ts.logger.Error("Authorization Error", "msg", err.Error(), "team", req)
		err = errNotAuthorizedToUpdate
		return
	}

	err = req.Validate()
	if err != nil {
		ts.logger.Error("Invalid request for team update", "msg", err.Error(), "team", req)
		return
	}

	updatedTeam, err := ts.store.UpdateTeam(ctx, id, db.Team{
		Name:        req.Name,
		Description: req.Description,
		ShowcaseUrl: req.ShowcaseUrl,
	})

	response.Team = &updatedTeam
	return
}

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection) Service {
	return &teamService{
		store:      s,
		logger:     l,
		collection: c,
	}
}
