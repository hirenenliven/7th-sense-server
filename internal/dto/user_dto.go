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
func (r *UserRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	if r.Name == "" {
		return response.ValidUserName
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return response.ValidEmail
	}
	return nil
}
