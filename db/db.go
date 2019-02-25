package db

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	// "github.com/pkg/errors"
)

type ctxKey int

const (
	dbKey          ctxKey = 0
	defaultTimeout        = 1 * time.Second
)

type Storer interface {
	// Category
	CreateCategory(ctx context.Context, collection *mongo.Collection, category *Category) (err error)
	ListCategories(ctx context.Context, collection *mongo.Collection) (categories []*Category, err error)
	FindCategoryByID(ctx context.Context, id string) (category Category, err error)
	DeleteCategoryByID(ctx context.Context, id string) (err error)
	UpdateCategory(ctx context.Context, collection *mongo.Collection, filter *Category, category *Category) (err error)

	CreateTeam(ctx context.Context, collection *mongo.Collection, team *Team) (err error)
	ListTeams(ctx context.Context, collection *mongo.Collection) (teams []*Team, err error)
	// Events
	CreateEvent(ctx context.Context, collection *mongo.Collection, event *Event) (err error)
	ListEvents(ctx context.Context, collection *mongo.Collection) (events []*Event, err error)
	FindEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (event Event, err error)
	DeleteEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (err error)
	UpdateEvent(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection, event *Event) (err error)
	// TeamMember
	CreateTeamMember(ctx context.Context, collection *mongo.Collection, teamMember *TeamMember) (err error)
	ListTeamMember(ctx context.Context, teamID primitive.ObjectID, collection *mongo.Collection) (teamMembers []*TeamMember, err error)
	FindTeamMemberByID(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection) (teamMember TeamMember, err error)
	DeleteTeamMemberByID(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection) (err error)
	UpdateTeamMember(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection, teamMember *TeamMember) (err error)
}

type store struct {
	db *mongo.Database
}

// func newContext(ctx context.Context, tx *sqlx.Tx) context.Context {
// 	return context.WithValue(ctx, dbKey, tx)
// }

// func Transact(ctx context.Context, dbx *mongo.Database, opts *sql.TxOptions, txFunc func(context.Context) error) (err error) {
// tx, err := dbx.BeginTxx(ctx, opts)
// if err != nil {
// 	return
// }
// defer func() {
// 	if p := recover(); p != nil {
// 		switch p := p.(type) {
// 		case error:
// 			err = errors.WithStack(p)
// 		default:
// 			err = errors.Errorf("%s", p)
// 		}
// 	}
// 	if err != nil {
// 		e := tx.Rollback()
// 		if e != nil {
// 			err = errors.WithStack(e)
// 		}
// 		return
// 	}
// 	err = errors.WithStack(tx.Commit())
// }()

// ctxWithTx := newContext(ctx, tx)
// err = WithDefaultTimeout(ctxWithTx, txFunc)
// return err
// }

func WithTimeout(ctx context.Context, timeout time.Duration, op func(ctx context.Context) error) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return op(ctxWithTimeout)
}

func WithDefaultTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	return WithTimeout(ctx, defaultTimeout, op)
}

func NewStorer(d *mongo.Database) Storer {
	return &store{
		db: d,
	}
}
