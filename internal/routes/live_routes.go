package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// registerLiveRoutes registers the Socket.IO live data handshake (ADR 0009) on the protected
// router, so authenticate resolves the session before the WebSocket upgrade. A "log in as" tab
// reaches it as /api/as/<key>/v1/socket.io/ (middleware.RouteLoginAs).
func registerLiveRoutes(router *http.ServeMux, prefix string, liveHandler *handler.LiveHandler) {
	router.HandleFunc("GET "+prefix+"/socket.io/", liveHandler.Serve)
}
