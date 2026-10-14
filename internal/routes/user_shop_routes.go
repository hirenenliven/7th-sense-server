package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Every user reads their company's rows; only managers change them (see targetUserID).
func registerUserShopRoutes(router *http.ServeMux, prefix string, userShopHandler *handler.UserShopHandler) {
	router.HandleFunc("GET "+prefix+"/user-shop", userShopHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-shop/{id}", userShopHandler.GetByID)
	router.Handle("POST "+prefix+"/user-shop", managersOnly(userShopHandler.Create))
	router.Handle("PUT "+prefix+"/user-shop/{id}", managersOnly(userShopHandler.Update))
	router.Handle("DELETE "+prefix+"/user-shop/{id}", managersOnly(userShopHandler.Delete))
}
