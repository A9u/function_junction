package event

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/A9u/function_junction/db"
	"go.uber.org/zap"
)

type Service interface {
	list(ctx context.Context) (response listResponse, err error)
	create(ctx context.Context, req createRequest) (response createResponse, err error)
	findByID(ctx context.Context, eventID primitive.ObjectID) (response findByIDResponse, err error)
	deleteByID(ctx context.Context, eventID primitive.ObjectID) (err error)
	update(ctx context.Context, req updateRequest, eventID primitive.ObjectID) (response updateResponse, err error)
}

type eventService struct {
	store      db.Storer
	logger     *zap.SugaredLogger
	collection *mongo.Collection
}

func (es *eventService) list(ctx context.Context) (response listResponse, err error) {
	events, err := es.store.ListEvents(ctx, es.collection)
	if err == db.ErrEventNotExist {
		es.logger.Error("No events present", "err", err.Error())
		return response, errNoEvents
	}
	if err != nil {
		es.logger.Error("Error listing events", "err", err.Error())
		return
	}

	response.Events = events
	return
}

func (es *eventService) create(ctx context.Context, c createRequest) (response createResponse, err error) {
	err = c.Validate()
	if err != nil {
		es.logger.Errorw("Invalid request for event create", "msg", err.Error(), "event", c)
		return
	}

	event, err := es.store.CreateEvent(ctx, es.collection, &db.Event{
		Title: c.Title,
		Description: c.Description,
		StartDateTime: c.StartDateTime,
		EndDateTime: c.EndDateTime,
		IsShowcasable: c.IsShowcasable,
		IsIndividualEvent: c.IsIndividualEvent,
		CreatedBy: c.CreatedBy,
		MaxSize: c.MaxSize,
		MinSize: c.MinSize,
		IsPublished: c.IsPublished,
		Venue: c.Venue,
		RegisterBefore: c.RegisterBefore,
	})

	if err != nil {
		es.logger.Error("Error creating event", "err", err.Error())
		return
	}
	response.Event = event
	return
}

func (es *eventService) findByID(ctx context.Context, id primitive.ObjectID) (response findByIDResponse, err error) {
	event, err := es.store.FindEventByID(ctx, id, es.collection)
	if err != nil {
		es.logger.Error("Error finding Event - ", "err", err.Error(), "event_id", id)
		return
	}
	response.Event = event
	return
}

func (es *eventService) update(ctx context.Context, eu updateRequest, id primitive.ObjectID) (response updateResponse, err error) {
	event, err := es.store.UpdateEvent(ctx, id, es.collection, &db.Event{Title: eu.Title, Description: eu.Description,
		Venue: eu.Venue, IsPublished: eu.IsPublished})
	if err != nil {
		es.logger.Error("Error updating event", "err", err.Error(), "event", eu)
		return
	}
	response.Event = event
	return
}

func (es *eventService) deleteByID(ctx context.Context, id primitive.ObjectID) (err error) {
	fmt.Println("I was here in service")
	err = es.store.DeleteEventByID(ctx, id, es.collection)
	if err != nil {
		es.logger.Error("Error deleting Event - ", "err", err.Error(), "event_id", id)
		return
	}
	return
}

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection) Service {
	return &eventService{
		store:      s,
		logger:     l,
		collection: c,
	}
}
