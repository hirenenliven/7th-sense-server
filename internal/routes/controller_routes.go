package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in can read (the machine form lists them); only the superadmin can change.
func registerControllerRoutes(router *http.ServeMux, prefix string, controllerHandler *handler.ControllerHandler) {
	router.HandleFunc("GET "+prefix+"/controllers", controllerHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/controllers/{id}", controllerHandler.GetByID)
	router.Handle("POST "+prefix+"/controllers", superAdminOnly(controllerHandler.Create))
	router.Handle("PUT "+prefix+"/controllers/{id}", superAdminOnly(controllerHandler.Update))
	router.Handle("DELETE "+prefix+"/controllers/{id}", superAdminOnly(controllerHandler.Delete))
}
