package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserStateHandler handles HTTP requests for the user's states (/user-state).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserStateHandler struct {
	service *service.UserStateService
}

// NewUserStateHandler creates a UserStateHandler that uses userStateService.
func NewUserStateHandler(userStateService *service.UserStateService) *UserStateHandler {
	return &UserStateHandler{service: userStateService}
}

// GetAll handles GET /api/v1/user-state[?userCountryId=<id>][&userId=<id>].
//
// @Summary      Get user states
// @Description  Returns the user's states. A non-nil userCountryID keeps only that country's states.
// @Description
// @Description  Used by pages: States, Cities, Add / edit customer or user, Customer / user detail
// @Tags         Organization
// @ID           getUserStates
// @Produce      json
// @Security     cookieAuth
// @Param        userCountryId query int false "User country id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserState}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-state [get]
func (userStateHandler *UserStateHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
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

// GetByID handles GET /api/v1/user-state/{id}[?userId=<id>].
func (userStateHandler *UserStateHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	state, err := userStateHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, state)
}

// Create handles POST /api/v1/user-state[?userId=<id>].
//
// @Summary      Create user state
// @Description  Adds a state for userID under one of the user's countries.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: States
// @Tags         Organization
// @ID           createUserState
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserStateRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserState}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-state [post]
func (userStateHandler *UserStateHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
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

// Delete handles DELETE /api/v1/user-state/{id}[?userId=<id>].
//
// @Summary      Delete user state
// @Description  Soft-deletes the state. A state that still has cities is kept: its cities must be deleted first.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: States
// @Tags         Organization
// @ID           deleteUserState
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
// @Router       /v1/user-state/{id} [delete]
func (userStateHandler *UserStateHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userStateHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "State deleted")
}
