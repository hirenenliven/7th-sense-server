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
func (userLocationHandler *UserLocationHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
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
