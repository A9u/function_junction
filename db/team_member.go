package db

import (
	"context"
	"fmt"
	"time"

	"github.com/A9u/function_junction/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type TeamMember struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status    string             `json:"status"`
	InviteeID primitive.ObjectID `json:"-"`
	InviterID primitive.ObjectID `json:"-"`
	TeamID    primitive.ObjectID `json:"-"`
	EventID   primitive.ObjectID `json:"event_id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type TeamMemberInfo struct {
	TeamMember
	InviteeInfo UserInfo `json:"Invitee"`
	InviterInfo UserInfo `json:"Inviter"`
}

type InviterInfo struct {
	InviterID   primitive.ObjectID `json:"invitee_id"`
	InviterName string             `json:"invitee_name"`
}

func (s *store) CreateTeamMember(ctx context.Context, collection *mongo.Collection, teamMember *TeamMember) (createdTeamMember TeamMember, err error) {
	fmt.Println("teamMemner", teamMember)
	teamMember.CreatedAt = time.Now()
	fmt.Println("teamMember", teamMember)

	res, err := collection.InsertOne(ctx, teamMember)
	if err != nil {
		fmt.Println("Error in CreateTeamMember: ", err)
		return
	}
	id := res.InsertedID
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&teamMember)
	return *teamMember, err
}

func (s *store) ListTeamMember(ctx context.Context, teamID primitive.ObjectID, eventID primitive.ObjectID, collection *mongo.Collection, userCollection *mongo.Collection, eventCollection *mongo.Collection, teamCollection *mongo.Collection) (teamMembers []*TeamMemberInfo, err error) {
	var invitee, inviter UserInfo
	cur, err := collection.Find(ctx, bson.D{{"teamid", teamID}})
	if err != nil {
		fmt.Println("Error in find: ", err)
		return teamMembers, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem TeamMember
		err = cur.Decode(&elem)
		invitee, err = FindUserInfoByID(ctx, elem.InviteeID)
		if err != nil {
			fmt.Println("Invitee does not exist:", err)
			return
		}

		inviter, _ = FindUserInfoByID(ctx, elem.InviterID)
		teamMemberInfo := TeamMemberInfo{TeamMember: elem, InviteeInfo: invitee, InviterInfo: inviter}
		teamMembers = append(teamMembers, &teamMemberInfo)
	}
	if err = cur.Err(); err != nil {
		fmt.Println("Error in currsor: ", err)
		return
	}
	return teamMembers, err
}

func (s *store) FindTeamMemberByID(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection) (teamMember TeamMember, err error) {
	err = collection.FindOne(ctx, bson.D{{"_id", teamMemberID}}).Decode(&teamMember)

	if err != nil {
		fmt.Println("Error in FindTeamMemberByID: ", err)
		return
	}
	return teamMember, err
}

func (s *store) FindTeamId(ctx context.Context, collection *mongo.Collection, inviterId primitive.ObjectID, eventID primitive.ObjectID) (teamId primitive.ObjectID, err error) {
	var teamMember TeamMember
	err = collection.FindOne(ctx, bson.D{{"inviterid", inviterId}, {"eventid", eventID}}).Decode(&teamMember)

	if err != nil {
		fmt.Println("Error in FindTeamMemberByID: ", err)
		return
	}
	return teamMember.TeamID, err
}

// TODO: func foo(id, name, address string)
// combine same types
func (s *store) IsTeamComplete(ctx context.Context, collection *mongo.Collection, teamID primitive.ObjectID, eventID primitive.ObjectID) (result bool, err error) {
	count, err := collection.Count(ctx, bson.D{{"status", "accept"}, {"teamid", teamID}, {"eventid", eventID}})
	// TODO: first check if error is there or count is 0
	// then move to success case
	if err == nil {
		event, err := s.FindEventByID(ctx, eventID)
		if err == nil {
			if int64(event.MaxSize) == count {
				return true, nil
			}
		}
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

func (s *store) DeleteTeamMemberByID(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection) (err error) {
	_, err = collection.DeleteOne(ctx, bson.D{{"_id", teamMemberID}})
	if err != nil {
		fmt.Println("Error in DeleteTeamMemberByID: ", err)
		return
	}
	return err
}

func (s *store) UpdateTeamMember(ctx context.Context, id primitive.ObjectID, collection *mongo.Collection, teamMember *TeamMember) (updateTeamMember TeamMember, err error) {

	_, err = collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{
		{"status", teamMember.Status},
		{"teamid", teamMember.TeamID},
		{"inviterid", teamMember.InviterID},
		{"inviteeid", teamMember.InviteeID},
		{"eventid", teamMember.EventID},
		{"updatedat", time.Now()},
	},
	},
	})
	if err != nil {
		fmt.Println("Error During UpdateTeamMember: ", err)
		return
	}

	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&teamMember)
	if err != nil {
		fmt.Println("Error During UpdateTeamMember: ", err)
		return
	}

	return *teamMember, err
}

func (s *store) FindTeamMemberByInviteeIDEventID(ctx context.Context, inviteeID primitive.ObjectID, eventID primitive.ObjectID, collection *mongo.Collection) (teamMember *TeamMember, err error) {

	err = collection.FindOne(ctx, bson.D{{"invitee_id", inviteeID}, {"event_id", eventID}, {"status", "accepted"}}).Decode(&teamMember)
	if err != nil {
		fmt.Println("Error During Finding team member: ", err)
		return
	}

	return teamMember, err
}

func (s *store) IsAttendingEvent(ctx context.Context, eventID primitive.ObjectID) (is_present bool) {
	collection := app.GetCollection("team_members")
	currentUserID := ctx.Value("currentUser").(User).ID
	err := collection.FindOne(ctx, bson.D{{"event_id", eventID}, {"status", "accepted"}, {"inviter_id", currentUserID}})
	if err != nil {
		err = collection.FindOne(ctx, bson.D{{"event_id", eventID}, {"status", "accepted"}, {"invitee_id", currentUserID}})
		if err != nil {
			is_present = false
		} else {
			is_present = true
		}
	} else {
		is_present = true
	}
	return
}

func (s *store) NumberOfIndividualsAttendingEvent(ctx context.Context, eventID primitive.ObjectID) (count int) {
	collection := app.GetCollection("team_members")
	var countattendees int64
	countattendees, _ = collection.Count(ctx, bson.D{{"event_id", eventID}, {"status", "accepted"}})
	count = int(countattendees)
	return
}

// TODO: instead of invitersInfo []*InviterInfo, we can simply use []InviterInfo
func (s *store) FindListOfInviters(ctx context.Context, currentUser User, userCollection *mongo.Collection, collection *mongo.Collection, eventID primitive.ObjectID) (invitersInfo []*InviterInfo, err error) {
	var user User
	// var usera *User
	// err = userCollection.FindOne(ctx,  bson.D{{"email", "priyanka@joshsoftware.com"}}).Decode(&usera)
	// fmt.Println("userid", usera.ID)
	cur, err := collection.Find(ctx, bson.D{{"eventid", eventID}, {"inviteeid", currentUser.ID}, {"status", "invited"}})
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem TeamMember
		err = cur.Decode(&elem)
		err = userCollection.FindOne(ctx, bson.D{{"_id", elem.InviterID}}).Decode(&user)
		inviterInfo := InviterInfo{InviterID: user.ID, InviterName: user.Email}
		invitersInfo = append(invitersInfo, &inviterInfo)
	}
	return invitersInfo, err
}

func (s *store) FindTeamMemberByInviteeIDTeamID(ctx context.Context, inviteeID primitive.ObjectID, teamID primitive.ObjectID) (teamMember TeamMember, err error) {
	collection := app.GetCollection("team_members")
	err = collection.FindOne(ctx, bson.D{{"invitee_id", inviteeID}, {"team_id", teamID}, {"status", "accepted"}}).Decode(&teamMember)
	if err != nil {
		fmt.Println("Error During Finding team member: ", err)
		return
	}
	return
}

