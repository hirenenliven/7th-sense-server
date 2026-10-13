package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserSupervisorRoutes(router *http.ServeMux, prefix string, userSupervisorHandler *handler.UserSupervisorHandler) {
	router.HandleFunc("GET "+prefix+"/user-supervisor", userSupervisorHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-supervisor/{id}", userSupervisorHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-supervisor", userSupervisorHandler.Create)
	router.HandleFunc("PUT "+prefix+"/user-supervisor/{id}", userSupervisorHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/user-supervisor/{id}", userSupervisorHandler.Delete)
}
