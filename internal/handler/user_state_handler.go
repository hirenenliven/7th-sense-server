package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserStateHandler handles HTTP requests for the user's states (/user-state).
// Every request acts for the current user (middleware.Authenticate), never for
// a user id sent by the client.
type UserStateHandler struct {
	service *service.UserStateService
}

// NewUserStateHandler creates a UserStateHandler that uses userStateService.
func NewUserStateHandler(userStateService *service.UserStateService) *UserStateHandler {
	return &UserStateHandler{service: userStateService}
}

// GetAll handles GET /api/v1/user-state[?userCountryId=<id>].
func (userStateHandler *UserStateHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userCountryID, err := queryID(request, "userCountryId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	states, err := userStateHandler.service.GetAll(request.Context(), userID, userCountryID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, states)
}

// GetByID handles GET /api/v1/user-state/{id}.
func (userStateHandler *UserStateHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	state, err := userStateHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, state)
}

// Create handles POST /api/v1/user-state.
func (userStateHandler *UserStateHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var stateRequest dto.CreateUserStateRequest
	if err := decodeAndValidate(writer, request, &stateRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	state, err := userStateHandler.service.Create(request.Context(), userID, stateRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, state)
}

// Delete handles DELETE /api/v1/user-state/{id}.
func (userStateHandler *UserStateHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userStateHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "State deleted")
}
