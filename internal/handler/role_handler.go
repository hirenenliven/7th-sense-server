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
//
// @Summary      Get assignable roles
// @Description  Returns the roles the user may give to the users they create or edit, by id. It is empty for users who manage no users.
// @Description
// @Description  Used by pages: Add / edit customer or user
// @Tags         Users
// @ID           getAssignableRoles
// @Produce      json
// @Security     cookieAuth
// @Success      200 {object} response.Envelope{data=[]model.Role}
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/roles [get]
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
