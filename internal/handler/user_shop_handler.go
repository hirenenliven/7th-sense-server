package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserShopHandler handles HTTP requests for the user's shops (/user-shop).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserShopHandler struct {
	service *service.UserShopService
}

// NewUserShopHandler creates a UserShopHandler that uses userShopService.
func NewUserShopHandler(userShopService *service.UserShopService) *UserShopHandler {
	return &UserShopHandler{service: userShopService}
}

// GetAll handles GET /api/v1/user-shop[?userLocationId=<id>][&userId=<id>].
//
// @Summary      Get user shops
// @Description  Returns the user's shops. A non-nil userLocationID keeps only that location's shops.
// @Description
// @Description  Used by pages: Shops, Lines, Add / edit customer or user, Customer / user detail
// @Tags         Organization
// @ID           getUserShops
// @Produce      json
// @Security     cookieAuth
// @Param        userLocationId query int false "User location id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserShop}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-shop [get]
func (userShopHandler *UserShopHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
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

// GetByID handles GET /api/v1/user-shop/{id}[?userId=<id>].
func (userShopHandler *UserShopHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
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

// Create handles POST /api/v1/user-shop[?userId=<id>].
//
// @Summary      Create user shop
// @Description  Adds a shop for userID under one of the user's locations.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Shops
// @Tags         Organization
// @ID           createUserShop
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserShopRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserShop}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-shop [post]
func (userShopHandler *UserShopHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
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

// Update handles PUT /api/v1/user-shop/{id}[?userId=<id>].
//
// @Summary      Update user shop
// @Description  Renames the shop and moves it under another of the user's locations.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Shops
// @Tags         Organization
// @ID           updateUserShop
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.UpdateUserShopRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.UserShop}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-shop/{id} [put]
func (userShopHandler *UserShopHandler) Update(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
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

// Delete handles DELETE /api/v1/user-shop/{id}[?userId=<id>].
//
// @Summary      Delete user shop
// @Description  Soft-deletes the shop. A shop that still has lines is kept: its lines must be deleted first. A shop a user's machine access names is kept too.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Shops
// @Tags         Organization
// @ID           deleteUserShop
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-shop/{id} [delete]
func (userShopHandler *UserShopHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
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
