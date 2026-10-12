package dto

import (
	"net/mail"
	"regexp"
	"strings"

	"7th-sense-server/internal/response"
)

var (
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,50}$`)
	mobilePattern   = regexp.MustCompile(`^\+?[0-9]{7,15}$`)
)

// User statuses accepted by UpdateUserRequest.
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

// CreateUserRequest is the body for POST /users.
type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Mobile    string `json:"mobile"`
	RoleID    int64  `json:"roleId"`
	BranchID  *int64 `json:"branchId"`
	CompanyID *int64 `json:"companyId"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateUserRequest) Validate() error {
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)
	request.Mobile = strings.TrimSpace(request.Mobile)

	if err := validateProfile(request.Username, request.Email, request.FirstName, request.LastName, request.Mobile); err != nil {
		return err
	}
	// bcrypt only uses the first 72 bytes of a password.
	if len(request.Password) < 8 || len(request.Password) > 72 {
		return response.ValidPassword
	}
	// Role is mandatory
	if request.RoleID <= 0 {
		return response.ValidRole
	}
	return nil
}

// UpdateUserRequest is the body for PUT /users/{id}. An empty status keeps
// the current one.
type UpdateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Mobile    string `json:"mobile"`
	RoleID    int64  `json:"roleId"`
	BranchID  *int64 `json:"branchId"`
	CompanyID *int64 `json:"companyId"`
	Status    string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateUserRequest) Validate() error {
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)
	request.Mobile = strings.TrimSpace(request.Mobile)
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))

	if err := validateProfile(request.Username, request.Email, request.FirstName, request.LastName, request.Mobile); err != nil {
		return err
	}
	if request.Status != "" && request.Status != UserStatusActive && request.Status != UserStatusInactive {
		return response.ValidStatus
	}
	// Role is mandatory
	if request.RoleID <= 0 {
		return response.ValidRole
	}
	return nil
}

// validateProfile checks the fields shared by create and update.
func validateProfile(username, email, firstName, lastName, mobile string) error {
	if !usernamePattern.MatchString(username) {
		return response.ValidUserName
	}
	// ParseAddress also accepts "Name <a@b.com>", so require a bare address.
	if address, err := mail.ParseAddress(email); err != nil || address.Address != email || len(email) > 255 {
		return response.ValidEmail
	}
	if len(firstName) > 100 || len(lastName) > 100 {
		return response.ValidName
	}
	if mobile != "" && !mobilePattern.MatchString(mobile) {
		return response.ValidMobile
	}
	return nil
}
