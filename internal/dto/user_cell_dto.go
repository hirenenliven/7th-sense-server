package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userCellNameMaxLength matches user_cell.name VARCHAR(100).
const userCellNameMaxLength = 100

// CreateUserCellRequest is the body for POST /user-cell.
type CreateUserCellRequest struct {
	UserLineID int64  `json:"userLineId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserCellRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserCell(request.Name, request.UserLineID)
}

// UpdateUserCellRequest is the body for PUT /user-cell/{id}: the new name and line.
type UpdateUserCellRequest struct {
	UserLineID int64  `json:"userLineId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *UpdateUserCellRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserCell(request.Name, request.UserLineID)
}

// validateUserCell checks the cell fields: a name and the line it belongs to.
func validateUserCell(name string, userLineID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userCellNameMaxLength {
		return response.ValidName
	}
	if userLineID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid line")
	}
	return nil
}
