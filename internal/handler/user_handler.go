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
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, users)
}

// GetByID handles GET /api/v1/users/{id}.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, user)
}

// Create handles POST /api/v1/users.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		response.SendError(w, err)
		return
	}

	user, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusCreated, user)
}

// Update handles PUT /api/v1/users/{id}.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		response.SendError(w, err)
		return
	}

	user, err := h.service.Update(r.Context(), r.PathValue("id"), req)
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, user)
}

// Delete handles DELETE /api/v1/users/{id}.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), r.PathValue("id")); err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, "User deleted")
}
