package team_member

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"fmt"
	"github.com/A9u/function_junction/config"
	"github.com/A9u/function_junction/db"
	"github.com/A9u/function_junction/mailer"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Service interface {
	list(ctx context.Context, teamID primitive.ObjectID) (response listResponse, err error)
	create(ctx context.Context, req createRequest, teamID primitive.ObjectID) (msg string, err error)
	findByID(ctx context.Context, teamMemberID primitive.ObjectID) (response findByIDResponse, err error)
	deleteByID(ctx context.Context, teamMemberID primitive.ObjectID) (err error)
	update(ctx context.Context, req updateRequest, teamMemberID primitive.ObjectID) (err error)
}

type teamMemberService struct {
	store          db.Storer
	logger         *zap.SugaredLogger
	collection     *mongo.Collection
	teamCollection *mongo.Collection
	userCollection *mongo.Collection
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

func (tms *teamMemberService) create(ctx context.Context, tm createRequest, teamID primitive.ObjectID) (msg string, err error) {
	err = tm.Validate()
	if err != nil {
		tms.logger.Errorw("Invalid request for team member create", "msg", err.Error(), "team member", tm)
		return
	}

	team, err := tms.store.FindTeamByID(ctx, teamID, tms.teamCollection)
	if err != nil {
		tms.logger.Errorw("Invalid request for team member create", "msg", err.Error(), "team member", tm)
		return
	}

	currentUser := ctx.Value("currentUser").(db.User)

	userErrEmails := ""
	userEmails := ""

	for i := 0; i < len(tm.Emails); i++ {
		user, err := db.FindUserByEmail(ctx, tm.Emails[i], tms.userCollection)

		if err == nil {
			_, err := tms.store.FindTeamMemberByInviteeIDEventID(ctx, user.ID, team.EventID, tms.collection)

			if err != nil {
				err = tms.store.CreateTeamMember(ctx, tms.collection, &db.TeamMember{
					InviteeID: user.ID,
					Status:    "Invited",
					InviterID: currentUser.ID,
					TeamID:    teamID,
					EventID:   team.EventID,
				})

				if err == nil {
					fmt.Println(user.Email)
					userEmails += user.Email + ","
				}
			} else {
				userErrEmails += tm.Emails[i] + ","
			}
		} else {
			userErrEmails += tm.Emails[i] + ","

		}
	}

	if len(userEmails) > 0 {
		userEmails = strings.TrimRight(userEmails, ",")
		invitees := strings.Split(userEmails, ",")
		notifyTeamMembers(invitees, team, currentUser, team.EventID)
		msg = strconv.Itoa(len(invitees)) + " invitations sent successfully"
	}

	if len(userErrEmails) > 0 {
		msg += "Invitations not sent for " + userErrEmails
		//tms.logger.Errorw("Error creating team member for " + userErrEmails + "err")
	}

	return msg, nil
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

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection, t *mongo.Collection, u *mongo.Collection) Service {
	return &teamMemberService{
		store:          s,
		logger:         l,
		collection:     c,
		teamCollection: t,
		userCollection: u,
	}
}

func notifyTeamMembers(invitees []string, team *db.Team, currentUser db.User, eventID primitive.ObjectID) {
	mail := mailer.Email{}
	mail.From = "anusha@joshsoftware.com" //currentUser.Email
	mail.To = invitees
	fmt.Println(mail.To)
	mail.Subject = "Invitation to join " + team.Name
	mail.Body = "I have invited you to join my team <b>" + team.Name + "</b>." +
		"<p> Please click <a href=" + config.URL() + "events/" + eventID.String() + " > here </a>. to see more details. <p>"

	mail.Send()
}
