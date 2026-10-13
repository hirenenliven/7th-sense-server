package middleware

import (
	"context"
	"net/http"
)

// devUserID is the user every request acts as until login is built.
// TODO(auth): replace with the user id resolved from the request's auth token.
const devUserID int64 = 6

type userIDKey struct{}

// Authenticate puts the current user's id in the request context. Handlers read
// it with UserID; a user id is never taken from the request body.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), userIDKey{}, devUserID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// UserID returns the current user's id set by Authenticate. ok is false when
// the request was not authenticated.
func UserID(ctx context.Context) (userID int64, ok bool) {
	userID, ok = ctx.Value(userIDKey{}).(int64)
	return userID, ok && userID > 0
}
