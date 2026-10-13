package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserShopRoutes(router *http.ServeMux, prefix string, userShopHandler *handler.UserShopHandler) {
	router.HandleFunc("GET "+prefix+"/user-shop", userShopHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-shop/{id}", userShopHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-shop", userShopHandler.Create)
	router.HandleFunc("PUT "+prefix+"/user-shop/{id}", userShopHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/user-shop/{id}", userShopHandler.Delete)
}
