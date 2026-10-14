package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Every user reads their company's rows; only managers change them (see targetUserID).
func registerUserCellRoutes(router *http.ServeMux, prefix string, userCellHandler *handler.UserCellHandler) {
	router.HandleFunc("GET "+prefix+"/user-cell", userCellHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-cell/{id}", userCellHandler.GetByID)
	router.Handle("POST "+prefix+"/user-cell", managersOnly(userCellHandler.Create))
	router.Handle("PUT "+prefix+"/user-cell/{id}", managersOnly(userCellHandler.Update))
	router.Handle("DELETE "+prefix+"/user-cell/{id}", managersOnly(userCellHandler.Delete))
}
