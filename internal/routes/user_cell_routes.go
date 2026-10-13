package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserCellRoutes(router *http.ServeMux, prefix string, userCellHandler *handler.UserCellHandler) {
	router.HandleFunc("GET "+prefix+"/user-cell", userCellHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-cell/{id}", userCellHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-cell", userCellHandler.Create)
	router.HandleFunc("PUT "+prefix+"/user-cell/{id}", userCellHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/user-cell/{id}", userCellHandler.Delete)
}
