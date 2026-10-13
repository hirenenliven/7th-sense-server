package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserCountryHandler handles HTTP requests for the user's countries (/user-country).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserCountryHandler struct {
	service *service.UserCountryService
}

// NewUserCountryHandler creates a UserCountryHandler that uses userCountryService.
func NewUserCountryHandler(userCountryService *service.UserCountryService) *UserCountryHandler {
	return &UserCountryHandler{service: userCountryService}
}

// GetAll handles GET /api/v1/user-country[?userId=<id>].
func (userCountryHandler *UserCountryHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	countries, err := userCountryHandler.service.GetAll(request.Context(), userID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, countries)
}

// GetByID handles GET /api/v1/user-country/{id}[?userId=<id>].
func (userCountryHandler *UserCountryHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	country, err := userCountryHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, country)
}

// Create handles POST /api/v1/user-country[?userId=<id>].
func (userCountryHandler *UserCountryHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var countryRequest dto.CreateUserCountryRequest
	if err := decodeAndValidate(writer, request, &countryRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	country, err := userCountryHandler.service.Create(request.Context(), userID, countryRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, country)
}

// Delete handles DELETE /api/v1/user-country/{id}[?userId=<id>].
func (userCountryHandler *UserCountryHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userCountryHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Country deleted")
}
