package handler

import (
	"net/http"

	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// RoleHandler handles HTTP requests for roles.
type RoleHandler struct {
	service *service.RoleService
}

// NewRoleHandler creates a RoleHandler that uses roleService.
func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{service: roleService}
}

// GetAll handles GET /api/v1/roles: the roles the current user may give to the
// users they create or edit (the role dropdown of the user form).
func (roleHandler *RoleHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	roles, err := roleHandler.service.GetAssignable(request.Context(), actor)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, roles)
}
