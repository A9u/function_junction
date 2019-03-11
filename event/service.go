package event

import (
	"context"
	"fmt"
	"time"

	"github.com/A9u/function_junction/config"
	"github.com/A9u/function_junction/db"
	"github.com/A9u/function_junction/mailer"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.uber.org/zap"
)

type Service interface {
	list(ctx context.Context) (response listResponse, err error)
	create(ctx context.Context, req createRequest) (response eventResponse, err error)
	findByID(ctx context.Context, eventID primitive.ObjectID) (response eventResponse, err error)
	deleteByID(ctx context.Context, eventID primitive.ObjectID) (err error)
	update(ctx context.Context, req updateRequest, eventID primitive.ObjectID) (response eventResponse, err error)
}

type eventService struct {
	store      db.Storer
	logger     *zap.SugaredLogger
	collection *mongo.Collection
}

func (es *eventService) list(ctx context.Context) (response listResponse, err error) {
	events, err := es.store.ListEvents(ctx)
	if err == db.ErrEventNotExist {
		es.logger.Error("No events present", "err", err.Error())
		// TODO: do not manually return if you already have named returns
		// assign to err object and simply call return
		return response, errNoEvents
	}
	if err != nil {
		es.logger.Error("Error listing events", "err", err.Error())
		return
	}

	response.Events = events

	return
}

// TODO: variable name c is not readable, names can only be a single letter when scope is less
func (es *eventService) create(ctx context.Context, c createRequest) (response eventResponse, err error) {
	// TODO: we can call this method only `Validate`.
	err = c.EventValidate()
	if err != nil {
		es.logger.Error("Invalid request for event create", "msg", err.Error(), "event", c)
		return
	}

	// TODO: add a check if this value can be type asserted into type db.User
	currentUser := ctx.Value("currentUser").(db.User)
	// TODO: do not use camel case in variable names
	event_info, err := es.store.CreateEvent(ctx, db.Event{
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

	if event_info.IsPublished {
		// TODO: lets use a mailerService so that we can mock this while testing
		notifyAll(event_info, currentUser)
	}

	response.Event = event_info
	return
}

func (es *eventService) findByID(ctx context.Context, id primitive.ObjectID) (response eventResponse, err error) {
	event_info, err := es.store.FindEventByID(ctx, id)
	if err != nil {
		es.logger.Error("Error finding Event - ", "err", err.Error(), "event_id", id)
		return
	}
	response.Event = event_info
	return
}

func (es *eventService) update(ctx context.Context, eu updateRequest, id primitive.ObjectID) (response eventResponse, err error) {
	currentUser := ctx.Value("currentUser").(db.User)
	// TODO: error should be checked immediately
	oldEvent, err := es.store.FindEventByID(ctx, id)

	if oldEvent.CreatedBy != currentUser.ID {
		err = errNotAuthorizedToUpdate
	}

	if err != nil {
		es.logger.Error("Authorization Error", "msg", err.Error(), "event", eu)
		return
	}

	err = eu.EventValidate()
	if err != nil {
		es.logger.Error("Invalid request for event update", "msg", err.Error(), "event", eu)
		return
	}
	event_info, err := es.store.UpdateEvent(ctx, id, db.Event{
		Title:             eu.Title,
		Description:       eu.Description,
		Venue:             eu.Venue,
		IsPublished:       eu.IsPublished,
		MinSize:           eu.MinSize,
		MaxSize:           eu.MaxSize,
		StartDateTime:     eu.StartDateTime,
		EndDateTime:       eu.EndDateTime,
		IsIndividualEvent: eu.IsIndividualEvent,
		RegisterBefore:    eu.RegisterBefore,
		IsShowcasable:     eu.IsShowcasable,
	})

	notifyOthers(oldEvent, event_info, currentUser)
	response.Event = event_info
	return
}

func (es *eventService) deleteByID(ctx context.Context, id primitive.ObjectID) (err error) {
	fmt.Println("I was here in service")
	err = es.store.DeleteEventByID(ctx, id)
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

func notifyAll(event db.EventInfo, currentUser db.User) {
	mail := mailer.Email{}
	mail.From = currentUser.Email
	mail.To = []string{config.AllEmail()}
	mail.Subject = "New Event added - " + event.Title
	mail.Body = "A new event <b>" + event.Title + "</b> has been added. " +
		"<p> It is at " + event.Venue + " from " + event.StartDateTime.Format(time.ANSIC) + " to " +
		event.EndDateTime.Format(time.ANSIC) + ". </p>" +
		"<p> Please check the details <a href=" + config.URL() + "events/" + getEventIDString(event.ID) + " > here </a> <p>"

	mail.Send()
}

func notifyOthers(oldEvent db.EventInfo, newEvent db.EventInfo, currentUser db.User) {
	if !oldEvent.IsPublished && newEvent.IsPublished {
		notifyAll(newEvent, currentUser)
	} else if oldEvent.Venue != newEvent.Venue || oldEvent.StartDateTime != newEvent.StartDateTime || oldEvent.EndDateTime != newEvent.EndDateTime {
		notifyChange(newEvent, currentUser)
	}
}

func notifyChange(event db.EventInfo, currentUser db.User) {
	mail := mailer.Email{}
	mail.From = currentUser.Email
	mail.To = []string{config.AllEmail()}

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
