package team_member

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/A9u/function_junction/db"
	"go.uber.org/zap"
)

type Service interface {
	list(ctx context.Context, teamID primitive.ObjectID) (response listResponse, err error)
	create(ctx context.Context, req createRequest) (err error)
	findByID(ctx context.Context, teamMemberID primitive.ObjectID) (response findByIDResponse, err error)
	deleteByID(ctx context.Context, teamMemberID primitive.ObjectID) (err error)
	update(ctx context.Context, req updateRequest, teamMemberID primitive.ObjectID) (err error)
}

type teamMemberService struct {
	store      db.Storer
	logger     *zap.SugaredLogger
	collection *mongo.Collection
}

func (tms *teamMemberService) list(ctx context.Context, teamID primitive.ObjectID) (response listResponse, err error) {
	teamMembers, err := tms.store.ListTeamMember(ctx, teamID, tms.collection)
	if err == db.ErrTeamMemberNotExist {
		tms.logger.Error("No team members added", "err", err.Error())
		return response, errNoTeamMember
	}
	if err != nil {
		tms.logger.Error("Error listing team members", "err", err.Error())
		return
	}

	response.TeamMembers = teamMembers
	return
}

func (tms *teamMemberService) create(ctx context.Context, tm createRequest) (err error) {
	err = tm.Validate()
	if err != nil {
		tms.logger.Errorw("Invalid request for team member create", "msg", err.Error(), "team member", tm)
		return
	}

	err = tms.store.CreateTeamMember(ctx, tms.collection, &db.TeamMember{Name: tm.Name, Status: tm.Status, InviterID: tm.InviterID, TeamID: tm.TeamID, EventID: tm.EventID})
	if err != nil {
		tms.logger.Error("Error creating team member", "err", err.Error())
		return
	}
	return
}

func (tms *teamMemberService) update(ctx context.Context, tm updateRequest, id primitive.ObjectID) (err error) {
	err = tm.Validate()
	if err != nil {
		tms.logger.Error("Invalid Request for team member update", "err", err.Error(), "team member", tm)
		return
	}

	err = tms.store.UpdateTeamMember(ctx, id, tms.collection, &db.TeamMember{Name: tm.Name, Status: tm.Status, InviterID: tm.InviterID, TeamID: tm.TeamID, EventID: tm.EventID})
	if err != nil {
		tms.logger.Error("Error updating team member", "err", err.Error(), "team member", tm)
		return
	}

	return
}

func (tms *teamMemberService) findByID(ctx context.Context, id primitive.ObjectID) (response findByIDResponse, err error) {
	teamMember, err := tms.store.FindTeamMemberByID(ctx, id, tms.collection)
	if err != nil {
		tms.logger.Error("Error finding Team Member", "err", err.Error(), "teammember_id", id)
		return
	}

	response.TeamMember = teamMember
	return
}

func (tms *teamMemberService) deleteByID(ctx context.Context, id primitive.ObjectID) (err error) {
	err = tms.store.DeleteTeamMemberByID(ctx, id, tms.collection)
	if err != nil {
		tms.logger.Error("Error deleting Team Member", "err", err.Error(), "team_member_id", id)
		return
	}

	return
}

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection) Service {
	return &teamMemberService{
		store:      s,
		logger:     l,
		collection: c,
	}
}
