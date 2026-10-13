package dto

import (
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// userCountryNameMaxLength matches user_country.name VARCHAR(100).
const userCountryNameMaxLength = 100

// CreateUserCountryRequest is the body for POST /user-country. A country is only added
// and removed, so there is no update request.
type CreateUserCountryRequest struct {
	Name string `json:"name"`
}

// Validate trims the name and checks it is usable.
func (request *CreateUserCountryRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" || utf8.RuneCountInString(request.Name) > userCountryNameMaxLength {
		return response.ValidName
	}
	return nil
}
