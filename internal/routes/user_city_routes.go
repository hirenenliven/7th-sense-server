package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserCityRoutes(router *http.ServeMux, prefix string, userCityHandler *handler.UserCityHandler) {
	router.HandleFunc("GET "+prefix+"/user-city", userCityHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-city/{id}", userCityHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-city", userCityHandler.Create)
	router.HandleFunc("PUT "+prefix+"/user-city/{id}", userCityHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/user-city/{id}", userCityHandler.Delete)
}
