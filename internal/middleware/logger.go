package middleware

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"time"
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

// Unwrap lets http.ResponseController and http.MaxBytesReader reach the
// underlying ResponseWriter.
func (recorder *statusRecorder) Unwrap() http.ResponseWriter {
	return recorder.ResponseWriter
}

// Hijack hands the connection to a WebSocket upgrade (the Socket.IO live data, ADR 0009): its
// library checks for http.Hijacker directly instead of going through Unwrap.
func (recorder *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	recorder.status = http.StatusSwitchingProtocols
	return http.NewResponseController(recorder.ResponseWriter).Hijack()
}

// Logger logs the method, path, status and duration of every request.
func Logger(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		recorder := &statusRecorder{ResponseWriter: writer, status: http.StatusOK}

		nextHandler.ServeHTTP(recorder, request)

		slog.Info("request",
			"method", request.Method,
			"path", request.URL.Path,
			"status", recorder.status,
			"duration", time.Since(startTime),
		)
	})
}
