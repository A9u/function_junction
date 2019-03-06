package server

import (
	"context"
	"net/http"

	"github.com/A9u/function_junction/db"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	user_id := r.Header.Get("Authorization")
	id, _ := primitive.ObjectIDFromHex(user_id)
	user, _ := db.FindUserByID(ctx, id)
	ctx = context.WithValue(ctx, "currentUser", user)

	next(rw, r.WithContext(ctx))
	// do some stuff after
}
