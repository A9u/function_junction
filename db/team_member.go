package db

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type TeamMember struct {
	Name      string             `json:"name"`
	Status    string             `json:"status"`
	InviteeID primitive.ObjectID `json:"invitee_id"`
	InviterID primitive.ObjectID `json:"inviter_id"`
	TeamID    primitive.ObjectID `json:"team_id"`
	EventID   primitive.ObjectID `json:"event_id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func (s *store) CreateTeamMember(ctx context.Context, collection *mongo.Collection, teamMember *TeamMember) (err error) {
	teamMember.CreatedAt = time.Now()
	fmt.Println("teamMember", teamMember)

	_, err = collection.InsertOne(ctx, teamMember)
	if err != nil {
		fmt.Println("Error in CreateTeamMember: ", err)
		return
	}
	return err
}

func (s *store) ListTeamMember(ctx context.Context, teamID primitive.ObjectID, collection *mongo.Collection) (teamMembers []*TeamMember, err error) {
	cur, err := collection.Find(ctx, bson.D{{"teamid", teamID}})
	if err != nil {
		fmt.Println("Error in find: ", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem TeamMember
		err = cur.Decode(&elem)
		teamMembers = append(teamMembers, &elem)
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

func (s *store) DeleteTeamMemberByID(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection) (err error) {
	_, err = collection.DeleteOne(ctx, bson.D{{"_id", teamMemberID}})
	if err != nil {
		fmt.Println("Error in DeleteTeamMemberByID: ", err)
		return
	}
	return err
}

func (s *store) UpdateTeamMember(ctx context.Context, id primitive.ObjectID, collection *mongo.Collection, teamMember *TeamMember) (err error) {

	_, err = collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{
		{"name", teamMember.Name},
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
	return err
}

func (s *store) FindTeamMemberByInviteeIDEventID(ctx context.Context, inviteeID primitive.ObjectID, eventID primitive.ObjectID, collection *mongo.Collection) (teamMember *TeamMember, err error) {

	err = collection.FindOne(ctx, bson.D{{"invitee_id", inviteeID}, {"event_id", eventID}, {"status", "accepted"}}).Decode(&teamMember)
	if err != nil {
		fmt.Println("Error During Finding team member: ", err)
		return
	}

	return teamMember, err
}
