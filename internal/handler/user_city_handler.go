package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserCityHandler handles HTTP requests for the user's cities (/user-city).
// Every request acts for the current user (middleware.Authenticate), never for
// a user id sent by the client.
type UserCityHandler struct {
	service *service.UserCityService
}

// NewUserCityHandler creates a UserCityHandler that uses userCityService.
func NewUserCityHandler(userCityService *service.UserCityService) *UserCityHandler {
	return &UserCityHandler{service: userCityService}
}

// GetAll handles GET /api/v1/user-city[?userStateId=<id>].
func (userCityHandler *UserCityHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
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

// GetByID handles GET /api/v1/user-city/{id}.
func (userCityHandler *UserCityHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	city, err := userCityHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, city)
}

// Create handles POST /api/v1/user-city.
func (userCityHandler *UserCityHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
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

// Update handles PUT /api/v1/user-city/{id}.
func (userCityHandler *UserCityHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

// Delete handles DELETE /api/v1/user-city/{id}.
func (userCityHandler *UserCityHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userCityHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "City deleted")
}
