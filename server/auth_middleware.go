package server

import (
	"context"
	"fmt"
	"net/http"
)

func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("testing middleware")
	ctx := r.Context()
	user_id := r.Header.Get("Authorization")
	ctx = context.WithValue(ctx, "user_id", user_id)

	next(rw, r.WithContext(ctx))
	// do some stuff after
}
