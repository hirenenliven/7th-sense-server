package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userStateNameMaxLength matches user_state.name VARCHAR(100).
const userStateNameMaxLength = 100

// CreateUserStateRequest is the body for POST /user-state. A state is only added
// and removed, so there is no update request.
type CreateUserStateRequest struct {
	UserCountryID int64  `json:"userCountryId"`
	Name          string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserStateRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserState(request.Name, request.UserCountryID)
}

// validateUserState checks the state fields: a name and the country it belongs to.
func validateUserState(name string, userCountryID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userStateNameMaxLength {
		return response.ValidName
	}
	if userCountryID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid country")
	}
	return nil
}
