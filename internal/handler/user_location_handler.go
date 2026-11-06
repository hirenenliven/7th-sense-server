package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserLocationHandler handles HTTP requests for the user's locations (/user-location).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserLocationHandler struct {
	service *service.UserLocationService
}

// NewUserLocationHandler creates a UserLocationHandler that uses userLocationService.
func NewUserLocationHandler(userLocationService *service.UserLocationService) *UserLocationHandler {
	return &UserLocationHandler{service: userLocationService}
}

// GetAll handles GET /api/v1/user-location[?branchId=<id>][&userId=<id>].
//
// @Summary      Get user locations
// @Description  Returns the user's locations. A non-nil branchID keeps only that branch's locations.
// @Description
// @Description  Used by pages: Locations, Shops, Add / edit customer or user, Customer / user detail
// @Tags         Organization
// @ID           getUserLocations
// @Produce      json
// @Security     cookieAuth
// @Param        branchId query int false "Branch id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserLocation}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-location [get]
func (userLocationHandler *UserLocationHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	branchID, err := queryID(request, "branchId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	locations, err := userLocationHandler.service.GetAll(request.Context(), userID, branchID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, locations)
}

// GetByID handles GET /api/v1/user-location/{id}[?userId=<id>].
func (userLocationHandler *UserLocationHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	location, err := userLocationHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, location)
}

// Create handles POST /api/v1/user-location[?userId=<id>].
//
// @Summary      Create user location
// @Description  Adds a location for userID under a branch.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Locations
// @Tags         Organization
// @ID           createUserLocation
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserLocationRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserLocation}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-location [post]
func (userLocationHandler *UserLocationHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var locationRequest dto.CreateUserLocationRequest
	if err := decodeAndValidate(writer, request, &locationRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	location, err := userLocationHandler.service.Create(request.Context(), userID, locationRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, location)
}

// Update handles PUT /api/v1/user-location/{id}[?userId=<id>].
//
// @Summary      Update user location
// @Description  Renames the location and moves it under another branch.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Locations
// @Tags         Organization
// @ID           updateUserLocation
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.UpdateUserLocationRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.UserLocation}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-location/{id} [put]
func (userLocationHandler *UserLocationHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var locationRequest dto.UpdateUserLocationRequest
	if err := decodeAndValidate(writer, request, &locationRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	location, err := userLocationHandler.service.Update(request.Context(), userID, id, locationRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, location)
}

// Delete handles DELETE /api/v1/user-location/{id}[?userId=<id>].
//
// @Summary      Delete user location
// @Description  Soft-deletes the location. A location that still has shops is kept: its shops must be deleted first.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Locations
// @Tags         Organization
// @ID           deleteUserLocation
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
// @Router       /v1/user-location/{id} [delete]
func (userLocationHandler *UserLocationHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userLocationHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Location deleted")
}
