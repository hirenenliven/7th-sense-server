package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Every user reads their company's rows; only managers change them (see targetUserID).
func registerUserSupervisorRoutes(router *http.ServeMux, prefix string, userSupervisorHandler *handler.UserSupervisorHandler) {
	router.HandleFunc("GET "+prefix+"/user-supervisor", userSupervisorHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-supervisor/{id}", userSupervisorHandler.GetByID)
	router.Handle("POST "+prefix+"/user-supervisor", managersOnly(userSupervisorHandler.Create))
	router.Handle("PUT "+prefix+"/user-supervisor/{id}", managersOnly(userSupervisorHandler.Update))
	router.Handle("DELETE "+prefix+"/user-supervisor/{id}", managersOnly(userSupervisorHandler.Delete))
}
