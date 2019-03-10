package event

import (
	"context"
	"fmt"
	"github.com/A9u/function_junction/config"
	"github.com/A9u/function_junction/db"
	"github.com/A9u/function_junction/mailer"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
	"time"
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

	currentUser := ctx.Value("currentUser").(db.User)

	event, err := es.store.CreateEvent(ctx, es.collection, &db.Event{
		Title:             c.Title,
		Description:       c.Description,
		StartDateTime:     c.StartDateTime,
		EndDateTime:       c.EndDateTime,
		IsShowcasable:     c.IsShowcasable,
		IsIndividualEvent: c.IsIndividualEvent,
		MaxSize:           c.MaxSize,
		MinSize:           c.MinSize,
		IsPublished:       c.IsPublished,
		Venue:             c.Venue,
		RegisterBefore:    c.RegisterBefore,
		CreatedBy:         currentUser.ID,
	})

	if err != nil {
		es.logger.Error("Error creating event", "err", err.Error())
		return
	}

	if event.IsPublished {
		notifyAll(event, currentUser)
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
	oldEvent, err := es.store.FindEventByID(ctx, id, es.collection)

	event, err := es.store.UpdateEvent(ctx, id, es.collection, &db.Event{Title: eu.Title, Description: eu.Description,
		Venue: eu.Venue, IsPublished: eu.IsPublished})

	if err != nil {
		es.logger.Error("Error updating event", "err", err.Error(), "event", eu)
		return
	}

	currentUser := ctx.Value("currentUser").(db.User)

	notifyOthers(oldEvent, event, currentUser)
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

func notifyAll(event *db.Event, currentUser db.User) {
	mail := mailer.Email{}
	mail.From = currentUser.Email
	mail.To = []string{"all@joshsoftware.com"}
	mail.Subject = "New Event added - " + event.Title
	mail.Body = "A new event <b>" + event.Title + "</b> has been added. " +
		"<p> It is at " + event.Venue + " from " + event.StartDateTime.Format(time.ANSIC) + " to " +
		event.EndDateTime.Format(time.ANSIC) + ". </p>" +
		"<p> Please check the details <a href=" + config.URL() + "events/" + getEventIDString(event.ID) + " > here </a> <p>"

	mail.Send()
}

func notifyOthers(oldEvent db.Event, newEvent *db.Event, currentUser db.User) {
	if !oldEvent.IsPublished && newEvent.IsPublished {
		notifyAll(newEvent, currentUser)
	} else if oldEvent.Venue != newEvent.Venue || oldEvent.StartDateTime != newEvent.StartDateTime || oldEvent.EndDateTime != newEvent.EndDateTime {
		notifyChange(newEvent, currentUser)
	}
}

func notifyChange(event *db.Event, currentUser db.User) {
	mail := mailer.Email{}
	mail.From = currentUser.Email
	mail.To = []string{"all@joshsoftware.com"}

	mail.Subject = "Event - " + event.Title + " has been updated"

	mail.Body = "The event - <b>" + event.Title + "</b> has been updated." +
		"<p> It is now at " + event.Venue + " from " + event.StartDateTime.Format(time.ANSIC) + " to " +
		event.EndDateTime.Format(time.ANSIC) + ". </p>" +
		"<p> Please check the details <a href=" + config.URL() + "events/" + getEventIDString(event.ID) + " > here </a> <p>"

	mail.Send()
}

func getEventIDString(eventID primitive.ObjectID) string {
	return eventID.Hex()
}
