package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userSupervisorNameMaxLength matches user_supervisor.name VARCHAR(100).
const userSupervisorNameMaxLength = 100

// CreateUserSupervisorRequest is the body for POST /user-supervisor.
type CreateUserSupervisorRequest struct {
	UserCellID int64  `json:"userCellId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserSupervisorRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserSupervisor(request.Name, request.UserCellID)
}

// UpdateUserSupervisorRequest is the body for PUT /user-supervisor/{id}: the new name and cell.
type UpdateUserSupervisorRequest struct {
	UserCellID int64  `json:"userCellId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *UpdateUserSupervisorRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserSupervisor(request.Name, request.UserCellID)
}

// validateUserSupervisor checks the supervisor fields: a name and the cell it belongs to.
func validateUserSupervisor(name string, userCellID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userSupervisorNameMaxLength {
		return response.ValidName
	}
	if userCellID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid cell")
	}
	return nil
}
