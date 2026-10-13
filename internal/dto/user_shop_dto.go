package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userShopNameMaxLength matches user_shop.name VARCHAR(100).
const userShopNameMaxLength = 100

// CreateUserShopRequest is the body for POST /user-shop.
type CreateUserShopRequest struct {
	UserLocationID int64  `json:"userLocationId"`
	Name           string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *CreateUserShopRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserShop(request.Name, request.UserLocationID)
}

// UpdateUserShopRequest is the body for PUT /user-shop/{id}: the new name and location.
type UpdateUserShopRequest struct {
	UserLocationID int64  `json:"userLocationId"`
	Name           string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *UpdateUserShopRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	return validateUserShop(request.Name, request.UserLocationID)
}

// validateUserShop checks the shop fields: a name and the location it belongs to.
func validateUserShop(name string, userLocationID int64) error {
	if name == "" || utf8.RuneCountInString(name) > userShopNameMaxLength {
		return response.ValidName
	}
	if userLocationID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid location")
	}
	return nil
}
