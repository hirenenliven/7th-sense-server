package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Every logged-in user reads the roles they may give (none for users who manage no users).
func registerRoleRoutes(router *http.ServeMux, prefix string, roleHandler *handler.RoleHandler) {
	router.HandleFunc("GET "+prefix+"/roles", roleHandler.GetAll)
}
