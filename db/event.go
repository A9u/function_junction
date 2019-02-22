package db

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"time"
)


type Event struct {
	// ID        string    `db:"id"`
	Id                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	StartDateTime     time.Time     `json:"startDateTime"`
	EndDateTime       time.Time     `json:"endDateTime"`
	IsShowcasable     bool          `json:"isShowcasable"`
	IsIndividualEvent bool          `json:"isIndividualParticipation"`
	CreatedBy         primitive.ObjectID   `json:"createdBy"`
	MaxSize           int           `json:"maxSize"`
	MinSize           int           `json:"minSize"`
	IsPublished       bool          `json:"isPublished"`
	Venue             string        `json:"venue"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
	RegisterBefore    time.Time     `db:"registerBefore"`
}

func (s *store) CreateEvent(ctx context.Context, collection *mongo.Collection, event *Event) (err error) {
	event.CreatedAt = time.Now()
	_, err = collection.InsertOne(ctx, event)
	return err
}


func (s *store) ListEvents(ctx context.Context, collection *mongo.Collection) (events []*Event, err error) {
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Error in find: ", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem Event
		err = cur.Decode(&elem)
		events = append(events, &elem)
	}
	if err := cur.Err(); err != nil {
	}
	return events, err
}

func (s *store) FindEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (event Event, err error) {
	err = collection.FindOne(ctx, bson.D{{"_id", eventID}}).Decode(&event)
	return event, err
}

func (s *store) DeleteEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (err error) {
	_, err = collection.DeleteOne(ctx, bson.D{{"_id", eventID}})
	return err
}

func (s *store) UpdateEvent(ctx context.Context, id primitive.ObjectID, collection *mongo.Collection, event *Event) (err error) {
	_, err = collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set",
		bson.D{ { "title", event.Title },
				{"description", event.Description },
				{"isPublished", event.IsPublished },
				{"venue", event.Venue },
				{"updated_at", time.Now() }, }, },})
	return err
}
