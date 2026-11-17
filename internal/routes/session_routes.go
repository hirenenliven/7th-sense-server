package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// registerSessionRoutes registers session management: every logged-in user's own sessions, and
// any user's sessions for the superadmin. /sessions/others is registered before /sessions/{id}
// for readers; ServeMux prefers the literal segment either way.
func registerSessionRoutes(router *http.ServeMux, prefix string, sessionHandler *handler.SessionHandler) {
	router.HandleFunc("GET "+prefix+"/sessions", sessionHandler.GetMine)
	router.HandleFunc("DELETE "+prefix+"/sessions/others", sessionHandler.DeleteMyOthers)
	router.HandleFunc("DELETE "+prefix+"/sessions/{id}", sessionHandler.DeleteMine)
	router.Handle("GET "+prefix+"/admin/users/{id}/sessions", superAdminOnly(sessionHandler.GetByUser))
	router.Handle("DELETE "+prefix+"/admin/users/{id}/sessions", superAdminOnly(sessionHandler.DeleteByUser))
	router.Handle("DELETE "+prefix+"/admin/sessions/{id}", superAdminOnly(sessionHandler.DeleteByID))
}
