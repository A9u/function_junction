package db

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
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

	CreateTeam(ctx context.Context, collection *mongo.Collection, team *Team) (createdTeam *TeamInfo, err error)
	ListTeams(ctx context.Context, collection *mongo.Collection, eventID primitive.ObjectID) (teams []*TeamInfo, err error)
	FindTeamByID(ctx context.Context, teamID primitive.ObjectID, collection *mongo.Collection) (team *Team, err error)
	DeleteTeamByID(ctx context.Context, teamID primitive.ObjectID, collection *mongo.Collection) (err error)

	// Events
	CreateEvent(ctx context.Context, collection *mongo.Collection, event *Event) (created_event *EventInfo, err error)
	ListEvents(ctx context.Context, collection *mongo.Collection) (events []*EventInfo, err error)
	FindEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (event *EventInfo, err error)
	DeleteEventByID(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection) (err error)
	UpdateEvent(ctx context.Context, eventID primitive.ObjectID, collection *mongo.Collection, event *Event) (updated_event *EventInfo, err error)

	// TeamMember
	CreateTeamMember(ctx context.Context, collection *mongo.Collection, teamMember *TeamMember) (createdTeamMember TeamMember, err error)
	ListTeamMember(ctx context.Context, teamID primitive.ObjectID, eventID primitive.ObjectID, collection *mongo.Collection, userCollection *mongo.Collection, eventCollection *mongo.Collection, teamCollection *mongo.Collection) (teamMembers []*TeamMemberInfo, err error)
	FindTeamMemberByID(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection) (teamMember TeamMember, err error)
	DeleteTeamMemberByID(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection) (err error)
	UpdateTeamMember(ctx context.Context, teamMemberID primitive.ObjectID, collection *mongo.Collection, teamMember *TeamMember) (updatedTeamMember TeamMember, err error)
	DeleteAllTeamMembers(ctx context.Context, teamID primitive.ObjectID) (err error)

	FindListOfInviters(ctx context.Context, currentUser User, userCollection *mongo.Collection, collection *mongo.Collection, eventID primitive.ObjectID) (invitersInfo []*InviterInfo, err error)
	FindTeamMemberByInviteeIDEventID(ctx context.Context, inviteeID primitive.ObjectID, eventID primitive.ObjectID, collection *mongo.Collection) (teamMember *TeamMember, err error)
	// User
	//FindUserByID(userID primitive.ObjectID) (user User, err error)
	// FindUserByEmail(ctx context.Context, email string)(user User, err error)
	IsTeamComplete(ctx context.Context, collection *mongo.Collection, teamID primitive.ObjectID, eventID primitive.ObjectID) (result bool, err error)
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
