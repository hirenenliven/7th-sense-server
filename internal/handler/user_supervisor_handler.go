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
