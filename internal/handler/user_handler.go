package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a UserHandler that uses userService.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

// GetAll handles GET /api/v1/users.
func (userHandler *UserHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	users, err := userHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, users)
}

// GetByID handles GET /api/v1/users/{id}.
func (userHandler *UserHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := userHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// Create handles POST /api/v1/users.
func (userHandler *UserHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var userRequest dto.CreateUserRequest
	if err := decodeAndValidate(writer, request, &userRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := userHandler.service.Create(request.Context(), userRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, user)
}

// Update handles PUT /api/v1/users/{id}.
func (userHandler *UserHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var userRequest dto.UpdateUserRequest
	if err := decodeAndValidate(writer, request, &userRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := userHandler.service.Update(request.Context(), id, userRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// Delete handles DELETE /api/v1/users/{id}.
func (userHandler *UserHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := userHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "User deleted")
}
