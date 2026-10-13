package dto

import (
	"net/http"
	"strings"

	"7th-sense-server/internal/response"
)

// LoginRequest is the body for POST /auth/login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"` // keep the session after the browser closes
}

// Validate trims the username and checks both fields are present. The format
// is not checked further: a wrong username fails as invalid credentials.
func (request *LoginRequest) Validate() error {
	request.Username = strings.TrimSpace(request.Username)

	if request.Username == "" {
		return response.NewError(http.StatusBadRequest, "Please enter username")
	}
	if request.Password == "" {
		return response.NewError(http.StatusBadRequest, "Please enter password")
	}
	// Longer values can never match: usernames are at most 50 characters and
	// bcrypt only uses the first 72 bytes of a password.
	if len(request.Username) > 50 || len(request.Password) > 72 {
		return response.InvalidCredentials
	}
	return nil
}
