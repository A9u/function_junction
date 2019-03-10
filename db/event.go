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

type Event struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title             string        		`json:"title"`
	Description       string        		`json:"description"`
	StartDateTime     time.Time     		`json:"start_date_time"`
	EndDateTime       time.Time     		`json:"end_date_time"`
	IsShowcasable     bool          		`json:"is_showcasable"`
	IsIndividualEvent bool          		`json:"is_individual_participation"`
	CreatedBy         primitive.ObjectID    `json:"-"`
	MaxSize           int           		`json:"max_size"`
	MinSize           int           		`json:"min_size"`
	IsPublished       bool          		`json:"is_published"`
	Venue             string        		`json:"venue"`
	CreatedAt         time.Time     		`json:"created_at"`
	UpdatedAt         time.Time     		`json:"updated_at"`
	RegisterBefore    time.Time     		`json:"register_before"`
}

type EventInfo struct {
  *Event
  CreatorInfo 			UserInfo `json:"created_by"`
  NumberOfParticipants	int 		`json:"number_of_participants"`
  IsAttending			bool		`json:"is_attending"`
}

func (s *store) CreateEvent(ctx context.Context, collection *mongo.Collection, event *Event) (eventInfo *EventInfo, err error) {
	event.CreatedAt = time.Now()
	if event.IsIndividualEvent == true{
		event.MinSize = 0
		event.MaxSize = 0
	}
	res, err := collection.InsertOne(ctx, event)

	id := res.InsertedID
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&event)
	eventInfo = getEventInfo(s, ctx, event)
	return eventInfo, err
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
		eventInfo := getEventInfo(s, ctx, &elem)
    	eventsInfo = append(eventsInfo, eventInfo)
	}
	if err := cur.Err(); err != nil {
	}
	return eventsInfo, err
}

func (s *store) FindEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (eventInfo *EventInfo, err error) {
	var event *Event
	err = collection.FindOne(ctx, bson.D{{"_id", eventID}}).Decode(&event)
	eventInfo = getEventInfo(s, ctx, event)
	return eventInfo, err
}

func (s *store) FindEventByName(ctx context.Context, eventName string) (eventID primitive.ObjectID , err error) {
	collection := app.GetCollection("events")
	var event Event
	err = collection.FindOne(ctx, bson.D{{"title", eventName}}).Decode(&event)
	return event.ID, err
}

func (s *store) DeleteEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (err error) {
	_, err = collection.DeleteOne(ctx, bson.D{{"_id", eventID}})
	return err
}

func (s *store) UpdateEvent(ctx context.Context, id primitive.ObjectID, collection *mongo.Collection, event *Event) (eventInfo *EventInfo, err error) {
	_, err = collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set",
		bson.D{ { "title", event.Title },
				{ "description", event.Description },
				{ "is_published", event.IsPublished },
				{ "venue", event.Venue },
				{ "start_date_time", event.StartDateTime },
				{ "end_date_time", event.EndDateTime },
				{ "is_showcasable", event.IsShowcasable },
				{ "is_individual_participation", event.IsIndividualEvent },
				{ "max_size", event.MaxSize },
				{ "min_size", event.MinSize },
				{ "register_before", event.RegisterBefore },
				{ "updated_at", time.Now() }, }, },})
	err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&event)
	eventInfo = getEventInfo(s, ctx, event)
	return eventInfo, err
}



func getEventInfo(s *store, ctx context.Context, event *Event) (eventInfo *EventInfo){
    creatorInfo, _ := FindUserInfoByID(ctx, event.CreatedBy)
    participants := 0
    if event.IsIndividualEvent{
    	participants = s.NumberOfIndividualsAttendingEvent(ctx, event.ID)
    } else {
    	teams, _ := s.ListTeams(ctx, app.GetCollection("teams"), event.ID)
    	participants = len(teams)
    }
    isAttending := s.IsAttendingEvent(ctx, event.ID)
	eventI := EventInfo{Event: event, CreatorInfo: creatorInfo, NumberOfParticipants: participants, IsAttending: isAttending}
	eventInfo = &eventI
	return eventInfo
}