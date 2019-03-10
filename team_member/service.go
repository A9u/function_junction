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
	list(ctx context.Context, teamID  primitive.ObjectID) (response listResponse, err error)
	create(ctx context.Context, req createRequest, teamID primitive.ObjectID) (response createResponse, err error)
	findByID(ctx context.Context, teamMemberID primitive.ObjectID) (response findByIDResponse, err error)
	deleteByID(ctx context.Context, teamMemberID primitive.ObjectID) (err error)
	update(ctx context.Context, req updateRequest, teamMemberID primitive.ObjectID, teamID primitive.ObjectID, eventID primitive.ObjectID) (response updateResponse, err error)
}

type teamMemberService struct {
	store          db.Storer
	logger         *zap.SugaredLogger
	collection     *mongo.Collection
	teamCollection *mongo.Collection
	userCollection *mongo.Collection
	eventCollection *mongo.Collection
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
	fmt.Println("team_members", teamMembers)
	response.TeamMembers = teamMembers
	return
}

func (tms *teamMemberService) create(ctx context.Context, tm createRequest, teamID primitive.ObjectID) (response createResponse, err error) {
	err = tm.Validate()
	if err != nil {
		tms.logger.Errorw("Invalid request for team member create", "msg", err.Error(), "team member", tm)
		return
	}

	team, err := tms.store.FindTeamByID(ctx, teamID, tms.teamCollection)
	fmt.Println("team", team)
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
				teamMember, err := tms.store.CreateTeamMember(ctx, tms.collection, &db.TeamMember{
					InviteeID: user.ID,
					Status:    "Invited",
					InviterID: currentUser.ID,
					TeamID:    teamID,
					EventID:   team.EventID,
				})
			response.TeamMember = teamMember
				
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
	var msg = ""
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
	response.Message = msg
	return 
}

func (tms *teamMemberService) update(ctx context.Context, tm updateRequest, id primitive.ObjectID, teamID primitive.ObjectID, eventID primitive.ObjectID) (response updateResponse, err error) {
	err = tm.Validate()
	if err != nil {
		tms.logger.Error("Invalid Request for team member update", "err", err.Error(), "team member", tm)
		return
	}
	// currentUser := ctx.Value("currentUser").(db.User)
	
	user, err := db.FindUserByID(ctx, tm.InviterID)
	fmt.Println("User", user)
	if err != nil{
		tms.logger.Error("Inviter Does not Exist in Db")
		err = errInviterDoesNotExist
		return
	}
	Event, err := tms.store.FindEventByID(ctx , eventID, tms.eventCollection)
	if err!= nil {
		tms.logger.Error("Event Does not Exist in Db")
		err = errEventDoesNotExist
		return
	}
	fmt.Println("event", Event.ID)

	Team,  err := tms.store.FindTeamByID(ctx, teamID, tms.teamCollection)
	if err!= nil {
		tms.logger.Error("Team Does not Exist in Db")
		err = errTeamDoesNotExist
		return
	}
	fmt.Println("event", Team.ID)

	if err != nil{
		tms.logger.Error("Team Not Present", err.Error())
		err = errTeamDoesNotExist
		return
	}
	if tm.Status == "accept"{
		result, _ :=   tms.store.IsTeamComplete(ctx, tms.collection, tm.TeamID, tm.EventID)
		if result ==  true {
			tms.logger.Error("Team is Already Complete")
			return 
		}
	}
	teamMember, err := tms.store.UpdateTeamMember(ctx, id, tms.collection, &db.TeamMember{ Status: tm.Status, InviterID: user.ID, TeamID: tm.TeamID, EventID: tm.EventID})
	if err != nil {
		tms.logger.Error("Error updating team member", "err", err.Error(), "team member", tm)
		return
	}
	// notifyTeamMembers(invitees, team, currentUser, team.EventID)
	// subject := currentUser.Email + " has " + tm.Status + " your Invite"
	// body := "Click on Link to see Event and Team details"
	// mailer.SendEmailTo("priyanka216668@gmail.com", user.Email, subject, body)
	//mailer.NotifyAll()
	
	response.TeamMember = teamMember
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

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection, t *mongo.Collection, u *mongo.Collection, e *mongo.Collection) Service {
	return &teamMemberService{
		store:          s,
		logger:         l,
		collection:     c,
		teamCollection: t,
		userCollection: u,
		eventCollection: e,
	}
}

func notifyTeamMembers(invitees []string, team *db.Team, currentUser db.User, eventID primitive.ObjectID) {
	mail := mailer.Email{}
	mail.From = "priyanka@joshsoftware.com"//currentUser.Email
	mail.To = invitees
	fmt.Println(mail.To)
	mail.Subject = "Invitation to join " + team.Name
	mail.Body = "I have invited you to join my team <b>" + team.Name + "</b>." +
		"<p> Please click <a href=" + config.URL() + "events/" + getStringID(eventID) + " > here </a>. to see more details. <p>"

	mail.Send()
}

func getStringID(id primitive.ObjectID) string {
	return id.Hex()
}
