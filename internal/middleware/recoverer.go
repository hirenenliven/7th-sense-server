package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"7th-sense-server/internal/response"
)

// Recoverer turns a panic in a handler into a 500 response instead of
// crashing the server.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v\n%s", err, debug.Stack())
				response.SendResponse(w, r, response.ServerError, nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
