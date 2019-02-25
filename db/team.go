package db

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Team struct {
	//ID        primitive.ObjectID `json:"_id"`
	Name        string             `json:"name"`
	EventID     primitive.ObjectID `json:"event_id"`
	CreatorID   primitive.ObjectID `json:"creator_id"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	ShowcaseUrl string             `json:"showcase_url"`
	Description string             `json:"description"`
}

func (s *store) CreateTeam(ctx context.Context, collection *mongo.Collection, team *Team) (err error) {
	now := time.Now()
	team.CreatedAt = now
	team.UpdatedAt = now
	_, err = collection.InsertOne(ctx, team)
	if err != nil {
		fmt.Println("Error in team creation ", err, team)
		return
	}
	return err
}

func (s *store) ListTeams(ctx context.Context, collection *mongo.Collection) (teams []*Team, err error) {
	// findOptions := options.Find()
	fmt.Println(collection)
	cur, err := collection.Find(ctx, bson.D{})
	fmt.Println(cur)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Error in find: ", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem Team
		err = cur.Decode(&elem)
		teams = append(teams, &elem)
	}
	if err := cur.Err(); err != nil {
	}
	return teams, err
}
