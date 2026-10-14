package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Every user reads their company's rows; only managers change them (see targetUserID).
func registerUserLineRoutes(router *http.ServeMux, prefix string, userLineHandler *handler.UserLineHandler) {
	router.HandleFunc("GET "+prefix+"/user-line", userLineHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-line/{id}", userLineHandler.GetByID)
	router.Handle("POST "+prefix+"/user-line", managersOnly(userLineHandler.Create))
	router.Handle("PUT "+prefix+"/user-line/{id}", managersOnly(userLineHandler.Update))
	router.Handle("DELETE "+prefix+"/user-line/{id}", managersOnly(userLineHandler.Delete))
}
