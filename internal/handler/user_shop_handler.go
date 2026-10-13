package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserShopHandler handles HTTP requests for the user's shops (/user-shop).
// Every request acts for the current user (middleware.Authenticate), never for
// a user id sent by the client.
type UserShopHandler struct {
	service *service.UserShopService
}

// NewUserShopHandler creates a UserShopHandler that uses userShopService.
func NewUserShopHandler(userShopService *service.UserShopService) *UserShopHandler {
	return &UserShopHandler{service: userShopService}
}

// GetAll handles GET /api/v1/user-shop[?userLocationId=<id>].
func (userShopHandler *UserShopHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userLocationID, err := queryID(request, "userLocationId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shops, err := userShopHandler.service.GetAll(request.Context(), userID, userLocationID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shops)
}

// GetByID handles GET /api/v1/user-shop/{id}.
func (userShopHandler *UserShopHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shop, err := userShopHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shop)
}

// Create handles POST /api/v1/user-shop.
func (userShopHandler *UserShopHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var shopRequest dto.CreateUserShopRequest
	if err := decodeAndValidate(writer, request, &shopRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shop, err := userShopHandler.service.Create(request.Context(), userID, shopRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, shop)
}

// Update handles PUT /api/v1/user-shop/{id}.
func (userShopHandler *UserShopHandler) Update(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var shopRequest dto.UpdateUserShopRequest
	if err := decodeAndValidate(writer, request, &shopRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shop, err := userShopHandler.service.Update(request.Context(), userID, id, shopRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shop)
}

// Delete handles DELETE /api/v1/user-shop/{id}.
func (userShopHandler *UserShopHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := userShopHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Shop deleted")
}
