package service

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// RoleService holds the business rules for roles: which roles a user may give
// (model.AssignableRoles).
type RoleService struct {
	roleRepository *repository.RoleRepository
}

// NewRoleService creates a RoleService backed by roleRepository.
func NewRoleService(roleRepository *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepository: roleRepository}
}

// GetAssignable returns the roles actor may give to the users they create or
// edit, by id. It is empty for users who manage no users.
func (roleService *RoleService) GetAssignable(ctx context.Context, actor model.AuthUser) ([]model.Role, error) {
	roles, err := roleService.roleRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(roles, func(role model.Role) bool {
		return !model.CanAssignRole(actor.Role, role.Name)
	}), nil
}

func (roleService *RoleService) GetByID(ctx context.Context, id int64) (model.Role, error) {
	role, err := roleService.roleRepository.FindByID(ctx, id)
	return role, toRoleResponseError(err)
}

// toRoleResponseError maps repository errors to role responses. Other errors
// pass through and are sent as ServerError.
func toRoleResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Role not found")
	default:
		return err
	}
}
