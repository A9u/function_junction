package db

import (
	"context"
	"fmt"
	"github.com/A9u/function_junction/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type TeamMember struct { 
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status    string             `json:"status"`
	InviteeID primitive.ObjectID `json:"invitee_id"`
	InviterID primitive.ObjectID `json:"inviter_id"`
	TeamID    primitive.ObjectID `json:"team_id"`
	EventID   primitive.ObjectID `json:"event_id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type TeamMemberInfo struct{
	TeamMember
	InviteeName string `json:"invitee_name"`
	InviterName string `json:"inviter_name"`
}
func (s *store) CreateTeamMember(ctx context.Context, collection *mongo.Collection, teamMember *TeamMember) (createdTeamMember TeamMember, err error) {
	fmt.Println("teamMemner" ,teamMember)
	teamMember.CreatedAt = time.Now()
	fmt.Println("teamMember", teamMember)

	res, err := collection.InsertOne(ctx, teamMember)
	if err != nil {
		fmt.Println("Error in CreateTeamMember: ", err)
		return
	}
	id :=  res.InsertedID
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&teamMember)
	return *teamMember, err
}

func (s *store) ListTeamMember(ctx context.Context, teamID primitive.ObjectID, collection *mongo.Collection) (teamMembers []*TeamMemberInfo, err error) {
	

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Error in find: ", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem TeamMember
		err = cur.Decode(&elem)
		teamMemberInfo := TeamMemberInfo{ TeamMember: elem, InviteeName: "test", InviterName: "inviter" }
		fmt.Println("info", teamMemberInfo)
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

func (s *store) FindTeamId(ctx context.Context, collection *mongo.Collection, inviterId  primitive.ObjectID, eventID primitive.ObjectID) (teamId  primitive.ObjectID, err error) {
	var teamMember TeamMember 
	err = collection.FindOne(ctx, bson.D{{"inviterid", inviterId}, {"eventid", eventID}}).Decode(&teamMember)

	if err != nil {
		fmt.Println("Error in FindTeamMemberByID: ", err)
		return
	}
	return teamMember.TeamID, err
}

func (s *store) IsTeamComplete(ctx context.Context, collection *mongo.Collection, teamID primitive.ObjectID, eventID primitive.ObjectID)( result bool, err error){
	count, err := collection.Count(ctx, bson.D{{"status", "accept"}, {"teamid", teamID}, {"eventid", eventID}})
	if err == nil {
		event, err := s.FindEventByID(ctx, eventID, app.GetCollection("events"))
		if err == nil {
			if (int64(event.MaxSize) == count){
			return  true, nil
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
		{"team_id", teamMember.TeamID},
		{"inviter_id", teamMember.InviterID},
		{"updated_at", time.Now()},
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
