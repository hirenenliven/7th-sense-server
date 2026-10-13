package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userLocationNameMaxLength matches user_location.name VARCHAR(100).
const userLocationNameMaxLength = 100

// CreateUserLocationRequest is the body for POST /user-location.
type CreateUserLocationRequest struct {
	BranchID int64  `json:"branchId"`
	Name     string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserLocationRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserLocation(request.Name, request.BranchID)
}

// UpdateUserLocationRequest is the body for PUT /user-location/{id}: the new name and branch.
type UpdateUserLocationRequest struct {
	BranchID int64  `json:"branchId"`
	Name     string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *UpdateUserLocationRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserLocation(request.Name, request.BranchID)
}

// validateUserLocation checks the location fields: a name and the branch it belongs to.
func validateUserLocation(name string, branchID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userLocationNameMaxLength {
		return response.ValidName
	}
	if branchID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid branch")
	}
	return nil
}
