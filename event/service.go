package event

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/A9u/function_junction/db"
	"go.uber.org/zap"
)

type Service interface {
	list(ctx context.Context) (response listResponse, err error)
	create(ctx context.Context, req createRequest) (err error)
	// findByID(ctx context.Context, id string) (response findByIDResponse, err error)
	// deleteByID(ctx context.Context, id string) (err error)
	// update(ctx context.Context, req updateRequest) (err error)
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

func (es *eventService) create(ctx context.Context, c createRequest) (err error) {
	err = c.Validate()
	if err != nil {
		es.logger.Errorw("Invalid request for event create", "msg", err.Error(), "event", c)
		return
	}

	err = es.store.CreateEvent(ctx, es.collection, &db.Event{
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
	})

	if err != nil {
		es.logger.Error("Error creating event", "err", err.Error())
		return
	}
	return
}

// func (es *eventService) update(ctx context.Context, c updateRequest) (err error) {
// 	err = c.Validate()
// 	if err != nil {
// 		es.logger.Error("Invalid Request for event update", "err", err.Error(), "event", c)
// 		return
// 	}
//
// 	err = es.store.UpdateEvent(ctx, es.collection, &db.Category{Name: c.Name}, &db.Category{Name: c.Set.Name, Type: c.Set.Type})
// 	if err != nil {
// 		es.logger.Error("Error updating event", "err", err.Error(), "event", c)
// 		return
// 	}
//
// 	return
// }
//
// func (es *eventService) findByID(ctx context.Context, id string) (response findByIDResponse, err error) {
// 	// category, err := es.store.FindCategoryByID(ctx, id)
// 	// if err == db.ErrCategoryNotExist {
// 	// 	es.logger.Error("No category present", "err", err.Error())
// 	// 	return response, errNoCategoryId
// 	// }
// 	// if err != nil {
// 	// 	es.logger.Error("Error finding category", "err", err.Error(), "category_id", id)
// 	// 	return
// 	// }
//
// 	// response.Category = category
// 	return
// }
//
// func (es *eventService) deleteByID(ctx context.Context, id string) (err error) {
// 	// err = es.store.DeleteCategoryByID(ctx, id)
// 	// if err == db.ErrCategoryNotExist {
// 	// 	es.logger.Error("Category Not present", "err", err.Error(), "category_id", id)
// 	// 	return errNoCategoryId
// 	// }
// 	// if err != nil {
// 	// 	es.logger.Error("Error deleting category", "err", err.Error(), "category_id", id)
// 	// 	return
// 	// }
//
// 	return
// }

func NewService(s db.Storer, l *zap.SugaredLogger, c *mongo.Collection) Service {
	return &eventService{
		store:      s,
		logger:     l,
		collection: c,
	}
}
