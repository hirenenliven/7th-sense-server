package dto

import (
	"net/mail"
	"strings"

	"7th-sense-server/internal/response"
)

// UserRequest is the body for creating or updating a user.
type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validate trims the fields and checks they are usable.
func (request *UserRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))

	if request.Name == "" {
		return response.ValidUserName
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return response.ValidEmail
	}
	return nil
}
