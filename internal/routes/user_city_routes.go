package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Every user reads their company's rows; only managers change them (see targetUserID).
func registerUserCityRoutes(router *http.ServeMux, prefix string, userCityHandler *handler.UserCityHandler) {
	router.HandleFunc("GET "+prefix+"/user-city", userCityHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-city/{id}", userCityHandler.GetByID)
	router.Handle("POST "+prefix+"/user-city", managersOnly(userCityHandler.Create))
	router.Handle("PUT "+prefix+"/user-city/{id}", managersOnly(userCityHandler.Update))
	router.Handle("DELETE "+prefix+"/user-city/{id}", managersOnly(userCityHandler.Delete))
}
