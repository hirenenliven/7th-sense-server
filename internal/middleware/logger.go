package middleware

import (
	"net/http"
	"time"

	"7th-sense-server/internal/logger"
)

// statusRecorder captures the status code written by the next handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (recorder *statusRecorder) WriteHeader(status int) {
	recorder.status = status
	recorder.ResponseWriter.WriteHeader(status)
}

// Logger logs the method, path, status and duration of every request.
func Logger(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		recorder := &statusRecorder{ResponseWriter: writer, status: http.StatusOK}

		nextHandler.ServeHTTP(recorder, request)

		logger.Log(logger.LevelInfo, "request",
			"method", request.Method,
			"path", request.URL.Path,
			"status", recorder.status,
			"duration", time.Since(startTime),
		)
	})
}
