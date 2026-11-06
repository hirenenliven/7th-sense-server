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
//
// @Summary      Get user countries
// @Description  Returns the user's countries.
// @Description
// @Description  Used by pages: Countries, States, Add / edit customer or user, Customer / user detail
// @Tags         Organization
// @ID           getUserCountries
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserCountry}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-country [get]
func (userCountryHandler *UserCountryHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
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
//
// @Summary      Create user country
// @Description  Adds a country for userID.
// @Description  Roles: superadmin only (others get 403).
// @Description  Validation: Please enter valid name.
// @Description
// @Description  Used by pages: Countries
// @Tags         Organization
// @ID           createUserCountry
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserCountryRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserCountry}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-country [post]
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
//
// @Summary      Delete user country
// @Description  Soft-deletes the country. A country that still has states is kept: its states must be deleted first.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Countries
// @Tags         Organization
// @ID           deleteUserCountry
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
// @Router       /v1/user-country/{id} [delete]
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
