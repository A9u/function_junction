package db

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Event struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	StartDateTime     time.Time          `json:"startDateTime"`
	EndDateTime       time.Time          `json:"endDateTime"`
	IsShowcasable     bool               `json:"isShowcasable"`
	IsIndividualEvent bool               `json:"isIndividualParticipation"`
	CreatedBy         primitive.ObjectID `json:"createdBy"`
	MaxSize           int                `json:"maxSize"`
	MinSize           int                `json:"minSize"`
	IsPublished       bool               `json:"isPublished"`
	Venue             string             `json:"venue"`
	CreatedAt         time.Time          `db:"createdAt"`
	UpdatedAt         time.Time          `db:"updatedAt"`
	RegisterBefore    time.Time          `db:"registerBefore"`
}

type EventInfo struct {
  Event
  UserFirstName         string
  UserLastName          string
}

func (s *store) CreateEvent(ctx context.Context, collection *mongo.Collection, event *Event) (created_event *Event, err error) {
	event.CreatedAt = time.Now()
	res, err := collection.InsertOne(ctx, event)
	//if err != nil { return res,err }
	id := res.InsertedID
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&event)
	return event, err
}

func (s *store) ListEvents(ctx context.Context, collection *mongo.Collection) (eventsInfo []*EventInfo, err error) {
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Error in find: ", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem Event
		err = cur.Decode(&elem)
    user, _ := FindUserByID(ctx, elem.CreatedBy)
    fmt.Println(user)
    event := EventInfo{Event: elem, UserFirstName: user.FirstName, UserLastName: user.LastName}
    eventsInfo = append(eventsInfo, &event)
	}
	if err := cur.Err(); err != nil {
	}

	return eventsInfo, err
}

func (s *store) FindEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (event Event, err error) {
	err = collection.FindOne(ctx, bson.D{{"_id", eventID}}).Decode(&event)
	return event, err
}

func (s *store) DeleteEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (err error) {
	_, err = collection.DeleteOne(ctx, bson.D{{"_id", eventID}})
	return err
}

func (s *store) UpdateEvent(ctx context.Context, id primitive.ObjectID, collection *mongo.Collection, event *Event) (updated_event *Event, err error) {
	event.UpdatedAt = time.Now()
	_, err = collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set",
		bson.D{{"title", event.Title},
			{"description", event.Description},
			{"isPublished", event.IsPublished},
			{"venue", event.Venue},
			{"startDateTime", event.StartDateTime},
			{"endDateTime", event.EndDateTime},
			{"isShowcasable", event.IsShowcasable},
			{"isIndividualParticipation", event.IsIndividualEvent},
			{"maxSize", event.MaxSize},
			{"minSize", event.MinSize},
			{"registerBefore", event.RegisterBefore},
			{"updated_at", time.Now()}}}})
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&event)
	return event, err
}
