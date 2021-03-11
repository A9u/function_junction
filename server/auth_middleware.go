package server

import (
	"context"
	"net/http"

	"github.com/joshsoftware/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	//TODO: we should check if user id is not present
	user_id := r.Header.Get("Authorization")
	//TODO: we should never ignore errors in Go
	id, _ := primitive.ObjectIDFromHex(user_id)
	user, _ := db.FindUserByID(ctx, id)
	ctx = context.WithValue(ctx, "currentUser", user)

	next(rw, r.WithContext(ctx))
	// do some stuff after
}
