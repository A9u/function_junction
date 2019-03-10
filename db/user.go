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
  FirstName  string         `json:"public_profile.first_name"`
	LastName  string             `json:"last_name"`
  PublicProfile interface{} `json:"public_profile" bson:"public_profile"`
	Email string             `json:"email"`
}

type publicProfile struct {
  FirstName string `json:"first_name" bson:"first_name"`
  LastName string  `json:"last_name" bson:"last_name"`
}

func name(publicProfileData interface{}) {
  //publicProfileRes, err := publicProfileData.(publicProfile)
}

func FindUserByID(ctx context.Context, userID primitive.ObjectID) (userDetails User, err error) {
	collection := app.GetCollection("users")
	err = collection.FindOne(ctx, bson.D{{"_id", userID}}).Decode(&userDetails)
  name(userDetails.PublicProfile)
	return userDetails, err
}

func FindUserByEmail(ctx context.Context, email string, collection *mongo.Collection) (userDetails User, err error) {
	err = collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&userDetails)

	return userDetails, err
}

// func (s *store) FindUserByEmail(ctx context.Context, email string) (user User, err error) {
// 	collection := app.GetCollection("users")
// 	err = collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
// 	return user, err
// }
