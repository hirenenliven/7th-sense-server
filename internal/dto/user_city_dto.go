package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userCityNameMaxLength matches user_city.name VARCHAR(100).
const userCityNameMaxLength = 100

// CreateUserCityRequest is the body for POST /user-city.
type CreateUserCityRequest struct {
	UserStateID int64  `json:"userStateId"`
	Name        string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserCityRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserCity(request.Name, request.UserStateID)
}

// UpdateUserCityRequest is the body for PUT /user-city/{id}: the new name and state.
type UpdateUserCityRequest struct {
	UserStateID int64  `json:"userStateId"`
	Name        string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *UpdateUserCityRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserCity(request.Name, request.UserStateID)
}

// validateUserCity checks the city fields: a name and the state it belongs to.
func validateUserCity(name string, userStateID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userCityNameMaxLength {
		return response.ValidName
	}
	if userStateID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid state")
	}
	return nil
}
