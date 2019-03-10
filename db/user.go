package db

import (
	"context"
	"github.com/A9u/function_junction/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type User struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
  FirstName  string         `bson:"first_name" json:"public_profile"`
	LastName  string             `json:"last_name"`
  PublicProfile publicProfile         `json: "public_profile" bson:"public_profile"`
	Email string             `json:"email"`
}

type publicProfile struct {
  FirstName string `json:"first_name"`
  LastName string  `json:"last_name"`
}
func FindUserByID(ctx context.Context, userID primitive.ObjectID) (userDetails User, err error) {
	collection := app.GetCollection("users")
	err = collection.FindOne(ctx, bson.D{{"_id", userID}}).Decode(&userDetails)
  userDetails.FirstName = "test"
  userDetails.LastName = "test"
	return userDetails, err
}

func FindUserByEmail(ctx context.Context, email string, collection *mongo.Collection) (userDetails User, err error) {
	err = collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&userDetails)

	return userDetails, err
}
