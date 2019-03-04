package db

import (
	"context"
	"github.com/A9u/function_junction/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type User struct {
	Id                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string        `json:"name"`
	Email             string        `json:"email"`
}

func (s *store) FindUserByID(ctx context.Context, userID primitive.ObjectID) (userDetails User, err error) {
	collection := app.GetCollection("users")
	err = collection.FindOne(ctx, bson.D{{"_id", userID}}).Decode(&userDetails)
	return userDetails, err
}
