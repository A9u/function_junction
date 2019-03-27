package team_member

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"fmt"

	"github.com/A9u/function_junction/config"
	"github.com/A9u/function_junction/constant"

	"github.com/A9u/function_junction/app"
	"github.com/A9u/function_junction/db"
	"github.com/A9u/function_junction/mailer"
	"go.uber.org/zap"
)

type Service interface {
	list(ctx context.Context, teamID primitive.ObjectID, eventID primitive.ObjectID) (response listResponse, err error)
	create(ctx context.Context, req createRequest, teamID primitive.ObjectID, eventID primitive.ObjectID) (response createResponse, err error)
	findByID(ctx context.Context, teamMemberID primitive.ObjectID) (response findByIDResponse, err error)
	deleteByID(ctx context.Context, teamMemberID primitive.ObjectID) (err error)
	update(ctx context.Context, req updateRequest, teamMemberID primitive.ObjectID, teamID primitive.ObjectID, eventID primitive.ObjectID) (response updateResponse, err error)
	findListOfInviters(ctx context.Context, eventID primitive.ObjectID) (response InviterslistResponse, err error)
}

type teamMemberService struct {
	store           db.Storer
	logger          *zap.SugaredLogger
	collection      *mongo.Collection
	teamCollection  *mongo.Collection
	userCollection  *mongo.Collection
	eventCollection *mongo.Collection
	mailer          mailer.MailerService
}

func (tms *teamMemberService) list(ctx context.Context, teamID primitive.ObjectID, eventID primitive.ObjectID) (response listResponse, err error) {
	teamMembers, err := tms.store.ListTeamMember(ctx, teamID, eventID, tms.collection, tms.userCollection, tms.eventCollection, tms.teamCollection)
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

func (tms *teamMemberService) findListOfInviters(ctx context.Context, eventID primitive.ObjectID) (response InviterslistResponse, err error) {
	currentUser := ctx.Value("currentUser").(db.User)
	invitersInfo, err := tms.store.FindListOfInviters(ctx, currentUser, tms.userCollection, tms.collection, eventID)
	// if err == db.ErrTeamMemberNotExist {
	// 	tms.logger.Error("No team members added", "err", err.Error())
	// 	return response, errNoTeamMember
	// }
	if err != nil {
		tms.logger.Error("Error listing Inviters Info", "err", err.Error())
		return
	}
	// TODO: remove extra prints
	fmt.Println("team_members", invitersInfo)
	response.InvitersInfo = invitersInfo
	return
}

func (tms *teamMemberService) create(ctx context.Context, tm createRequest, teamID primitive.ObjectID, eventID primitive.ObjectID) (response createResponse, err error) {

	currentUser := ctx.Value("currentUser").(db.User)
	zeroValue, _ := primitive.ObjectIDFromHex("")
	event, _ := tms.store.FindEventByID(ctx, eventID)
	fmt.Println(event)
	team, _ := tms.store.FindTeamByEventIDAndName(ctx, eventID, event.Title, app.GetCollection("teams"))
	if teamID == zeroValue {

		_, err = tms.store.CreateTeamMember(ctx, tms.collection, &db.TeamMember{
			InviteeID: currentUser.ID,
			Status:    constant.Accepted,
			TeamID:    team.ID,
			EventID:   team.EventID,
		})
		if err != nil {
			return
		}
		response.Message = "RSVP Done!"
		return
	} else {
		err = tm.Validate()
		if err != nil {
			tms.logger.Errorw("Invalid request for team member create", "msg", err.Error(), "team member", tm)
			return
		}

		// TODO: assign empty variables like: var foo string
		_, err = tms.store.FindTeamMemberByInviteeIDEventID(ctx, currentUser.ID, eventID, tms.collection)

		if err != nil {
			tms.logger.Errorw("Only accepted members can invite", "msg", err.Error(), "team member", tm)
			return
		}

		var userEmails, userErrEmails []string

		// TODO: use range and remove branching
		for _, email := range tm.Emails {
			user, err := db.FindUserByEmail(ctx, email, tms.userCollection)

			if err != nil {
				userErrEmails = append(userErrEmails, email)
				continue
			}

			_, err = tms.store.FindTeamMemberByInviteeIDEventID(ctx, user.ID, team.EventID, tms.collection)

			if err == nil {
				userErrEmails = append(userErrEmails, email)
				continue
			}

			_, err = tms.store.FindInvitedTeamMember(ctx, team.ID, user.ID, tms.collection)

			if err != nil {
				_, err = tms.store.CreateTeamMember(ctx, tms.collection, &db.TeamMember{
					InviteeID: user.ID,
					Status:    constant.Invited,
					InviterID: currentUser.ID,
					TeamID:    teamID,
					EventID:   team.EventID,
				})
			}

			if err != nil {
				userErrEmails = append(userErrEmails, email)
				continue
			}

			userEmails = append(userEmails, email)
		}

		if len(userEmails) > 0 {
			tms.notifyTeamMembers(userEmails, team, currentUser, team.EventID)
		}

		response.FailedEmails = userErrEmails
		return
	}
}

func (tms *teamMemberService) update(ctx context.Context, tm updateRequest, id primitive.ObjectID, teamID primitive.ObjectID, eventID primitive.ObjectID) (response updateResponse, err error) {
	err = tm.Validate()
	if err != nil {
		tms.logger.Error("Invalid Request for team member update", "err", err.Error(), "team member", tm)
		return
	}

	currentUser := ctx.Value("currentUser").(db.User)

	teamMember, err := tms.store.FindTeamMemberByID(ctx, id, tms.collection)
	fmt.Println("teamMember", teamMember)
	if err != nil {
		tms.logger.Error("Team Member Does not Exist in Db")
		err = errTeamMemberDoesNotExist
		return
	}

	event, err := tms.store.FindEventByID(ctx, eventID)
	if err != nil {
		tms.logger.Error("Event Does not Exist in Db")
		err = errEventDoesNotExist
		return
	}
	fmt.Println("event", event.ID)

	team, err := tms.store.FindTeamByID(ctx, teamID, tms.teamCollection)
	if err != nil {
		tms.logger.Error("Team Does not Exist in Db")
		err = errTeamDoesNotExist
		return
	}
	fmt.Println("team", team.ID)

	if err != nil {
		tms.logger.Error("Team Not Present", err.Error())
		err = errTeamDoesNotExist
		return
	}
	if tm.Status == constant.Accepted {
		result, _ := tms.store.IsTeamComplete(ctx, tms.collection, teamID, eventID)
		if result == true {
			tms.logger.Error("Team is Already Complete")
			return
		}
	}
	teamMemberInfo, err := tms.store.UpdateTeamMember(ctx, id, tms.collection, &db.TeamMember{Status: tm.Status, InviterID: teamMember.InviterID, InviteeID: teamMember.InviteeID, TeamID: teamID, EventID: eventID})
	if err != nil {
		tms.logger.Error("Error updating team member", "err", err.Error(), "team member", tm)
		return
	}

	// inviter, err := db.FindUserByID(ctx, teamMember.InviterID)
	inviter := db.User{}
	tms.notifyTeamMemberInvitationStatus(inviter, currentUser, team, teamMember)

	response.TeamMember = teamMemberInfo
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

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection, t *mongo.Collection, u *mongo.Collection, e *mongo.Collection, m mailer.MailerService) Service {
	return &teamMemberService{
		store:           s,
		logger:          l,
		collection:      c,
		teamCollection:  t,
		userCollection:  u,
		eventCollection: e,
		mailer:          m,
	}
}

func (tms *teamMemberService) notifyTeamMembers(invitees []string, team *db.Team, currentUser db.User, eventID primitive.ObjectID) {
	body := "I have invited you to join my team <b>" + team.Name + "</b>." +
		"<p> Please click <a href=" + config.URL() + "events/" + getStringID(eventID) + " > here </a>. to see more details. <p>"

	tms.mailer.Send(invitees, currentUser.Email, "Invitation to join "+team.Name, body)
}

func (tms *teamMemberService) notifyTeamMemberInvitationStatus(inviter db.User, invitee db.User, team *db.Team, teamMember db.TeamMember) {
	body := "I have " + teamMember.Status + " your invitation  to join  team <b>" + team.Name + "</b>."
	tms.mailer.Send([]string{inviter.Email}, invitee.Email, "Invitation "+teamMember.Status+"By"+invitee.Email, body)
}

func getStringID(id primitive.ObjectID) string {
	return id.Hex()
}
