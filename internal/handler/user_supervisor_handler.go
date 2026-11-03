package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserSupervisorHandler handles HTTP requests for the user's supervisors (/user-supervisor).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserSupervisorHandler struct {
	service *service.UserSupervisorService
}

// NewUserSupervisorHandler creates a UserSupervisorHandler that uses userSupervisorService.
func NewUserSupervisorHandler(userSupervisorService *service.UserSupervisorService) *UserSupervisorHandler {
	return &UserSupervisorHandler{service: userSupervisorService}
}

// GetAll handles GET /api/v1/user-supervisor[?userCellId=<id>][&userId=<id>].
//
// @Summary      Get user supervisors
// @Description  Returns the user's supervisors. A non-nil userCellID keeps only that cell's supervisors.
// @Description
// @Description  Used by pages: Supervisors
// @Tags         Organization
// @ID           getUserSupervisors
// @Produce      json
// @Security     cookieAuth
// @Param        userCellId query int false "User cell id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserSupervisor}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-supervisor [get]
func (userSupervisorHandler *UserSupervisorHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userCellID, err := queryID(request, "userCellId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	supervisors, err := userSupervisorHandler.service.GetAll(request.Context(), userID, userCellID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, supervisors)
}

// GetByID handles GET /api/v1/user-supervisor/{id}[?userId=<id>].
func (userSupervisorHandler *UserSupervisorHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	supervisor, err := userSupervisorHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, supervisor)
}

// Create handles POST /api/v1/user-supervisor[?userId=<id>].
//
// @Summary      Create user supervisor
// @Description  Adds a supervisor for userID under one of the user's cells.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Supervisors
// @Tags         Organization
// @ID           createUserSupervisor
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserSupervisorRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserSupervisor}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-supervisor [post]
func (userSupervisorHandler *UserSupervisorHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var supervisorRequest dto.CreateUserSupervisorRequest
	if err := decodeAndValidate(writer, request, &supervisorRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	supervisor, err := userSupervisorHandler.service.Create(request.Context(), userID, supervisorRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, supervisor)
}

// Update handles PUT /api/v1/user-supervisor/{id}[?userId=<id>].
//
// @Summary      Update user supervisor
// @Description  Renames the supervisor and moves it under another of the user's cells.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Supervisors
// @Tags         Organization
// @ID           updateUserSupervisor
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.UpdateUserSupervisorRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.UserSupervisor}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-supervisor/{id} [put]
func (userSupervisorHandler *UserSupervisorHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var supervisorRequest dto.UpdateUserSupervisorRequest
	if err := decodeAndValidate(writer, request, &supervisorRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	supervisor, err := userSupervisorHandler.service.Update(request.Context(), userID, id, supervisorRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, supervisor)
}

// Delete handles DELETE /api/v1/user-supervisor/{id}[?userId=<id>].
//
// @Summary      Delete user supervisor
// @Description  Soft-deletes the supervisor.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Supervisors
// @Tags         Organization
// @ID           deleteUserSupervisor
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
// @Router       /v1/user-supervisor/{id} [delete]
func (userSupervisorHandler *UserSupervisorHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userSupervisorHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Supervisor deleted")
}
