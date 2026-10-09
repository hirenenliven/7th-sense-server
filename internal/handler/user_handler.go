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

// NewUserHandler creates a UserHandler that uses svc.
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// GetAll handles GET /api/v1/users.
func (h *UserHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	users, err := h.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, users)
}

// GetByID handles GET /api/v1/users/{id}.
func (h *UserHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	user, err := h.service.GetByID(request.Context(), request.PathValue("id"))
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// Create handles POST /api/v1/users.
func (h *UserHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var req dto.UserRequest
	if err := decodeAndValidate(writer, request, &req); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := h.service.Create(request.Context(), req)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, user)
}

// Update handles PUT /api/v1/users/{id}.
func (h *UserHandler) Update(writer http.ResponseWriter, request *http.Request) {
	var req dto.UserRequest
	if err := decodeAndValidate(writer, request, &req); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := h.service.Update(request.Context(), request.PathValue("id"), req)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// Delete handles DELETE /api/v1/users/{id}.
func (h *UserHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	if err := h.service.Delete(request.Context(), request.PathValue("id")); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "User deleted")
}
