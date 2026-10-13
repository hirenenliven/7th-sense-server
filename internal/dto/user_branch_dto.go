package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userBranchNameMaxLength matches user_branch.name VARCHAR(100).
const userBranchNameMaxLength = 100

// CreateUserBranchRequest is the body for POST /user-branch.
type CreateUserBranchRequest struct {
	UserCityID int64  `json:"userCityId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserBranchRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserBranch(request.Name, request.UserCityID)
}

// UpdateUserBranchRequest is the body for PUT /user-branch/{id}: the new name and city.
type UpdateUserBranchRequest struct {
	UserCityID int64  `json:"userCityId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *UpdateUserBranchRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserBranch(request.Name, request.UserCityID)
}

// validateUserBranch checks the branch fields: a name and the city it belongs to.
func validateUserBranch(name string, userCityID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userBranchNameMaxLength {
		return response.ValidName
	}
	if userCityID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid city")
	}
	return nil
}
