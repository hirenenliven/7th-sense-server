package middleware

import (
	"net/http"
	"runtime/debug"

	"7th-sense-server/internal/logger"
	"7th-sense-server/internal/response"
)

// Recoverer turns a panic in a handler into a 500 response instead of
// crashing the server.
func Recoverer(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log(logger.LevelError, "panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)
				response.SendResponse(writer, request, response.ServerError, nil)
			}
		}()
		nextHandler.ServeHTTP(writer, request)
	})
}
