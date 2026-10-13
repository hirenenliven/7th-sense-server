package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userLineNameMaxLength matches user_line.name VARCHAR(100).
const userLineNameMaxLength = 100

// CreateUserLineRequest is the body for POST /user-line.
type CreateUserLineRequest struct {
	UserShopID int64  `json:"userShopId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserLineRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserLine(request.Name, request.UserShopID)
}

// UpdateUserLineRequest is the body for PUT /user-line/{id}: the new name and shop.
type UpdateUserLineRequest struct {
	UserShopID int64  `json:"userShopId"`
	Name       string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *UpdateUserLineRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserLine(request.Name, request.UserShopID)
}

// validateUserLine checks the line fields: a name and the shop it belongs to.
func validateUserLine(name string, userShopID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userLineNameMaxLength {
		return response.ValidName
	}
	if userShopID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid shop")
	}
	return nil
}
