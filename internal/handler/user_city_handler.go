package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserCityHandler handles HTTP requests for the user's cities (/user-city).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserCityHandler struct {
	service *service.UserCityService
}

// NewUserCityHandler creates a UserCityHandler that uses userCityService.
func NewUserCityHandler(userCityService *service.UserCityService) *UserCityHandler {
	return &UserCityHandler{service: userCityService}
}

// GetAll handles GET /api/v1/user-city[?userStateId=<id>][&userId=<id>].
//
// @Summary      Get user cities
// @Description  Returns the user's cities. A non-nil userStateID keeps only that state's cities.
// @Description
// @Description  Used by pages: Add / edit branch, Cities, Add / edit customer or user, Customer / user detail
// @Tags         Organization
// @ID           getUserCities
// @Produce      json
// @Security     cookieAuth
// @Param        userStateId query int false "User state id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserCity}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-city [get]
func (userCityHandler *UserCityHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userStateID, err := queryID(request, "userStateId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	cities, err := userCityHandler.service.GetAll(request.Context(), userID, userStateID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, cities)
}

// GetByID handles GET /api/v1/user-city/{id}[?userId=<id>].
func (userCityHandler *UserCityHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	city, err := userCityHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, city)
}

// Create handles POST /api/v1/user-city[?userId=<id>].
//
// @Summary      Create user city
// @Description  Adds a city for userID under one of the user's states.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Cities
// @Tags         Organization
// @ID           createUserCity
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserCityRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserCity}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-city [post]
func (userCityHandler *UserCityHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var cityRequest dto.CreateUserCityRequest
	if err := decodeAndValidate(writer, request, &cityRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	city, err := userCityHandler.service.Create(request.Context(), userID, cityRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, city)
}

// Update handles PUT /api/v1/user-city/{id}[?userId=<id>].
//
// @Summary      Update user city
// @Description  Renames the city and moves it under another of the user's states.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Cities
// @Tags         Organization
// @ID           updateUserCity
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.UpdateUserCityRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.UserCity}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-city/{id} [put]
func (userCityHandler *UserCityHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var cityRequest dto.UpdateUserCityRequest
	if err := decodeAndValidate(writer, request, &cityRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	city, err := userCityHandler.service.Update(request.Context(), userID, id, cityRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, city)
}

// Delete handles DELETE /api/v1/user-city/{id}[?userId=<id>].
//
// @Summary      Delete user city
// @Description  Soft-deletes the city. A city that still has branches is kept: its branches must be deleted first.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Cities
// @Tags         Organization
// @ID           deleteUserCity
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
// @Router       /v1/user-city/{id} [delete]
func (userCityHandler *UserCityHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userCityHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "City deleted")
}
