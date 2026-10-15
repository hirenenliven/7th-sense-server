package dto

import (
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"7th-sense-server/internal/response"
)

var (
	codePattern      = regexp.MustCompile(`^[A-Z0-9]{2,20}$`)
	gstPattern       = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[A-Z0-9]{1}[Z]{1}[A-Z0-9]{1}$`)
	panPattern       = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)
	colorCodePattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
)

// Branch statuses accepted by UpdateBranchRequest.
const (
	BranchStatusActive   = "active"
	BranchStatusInactive = "inactive"
)

// CreateBranchRequest is the body for POST /branches.
type CreateBranchRequest struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	UserCityID *int64 `json:"userCityId"`
	CustomerID *int64 `json:"customerId"`
	GSTNumber  string `json:"gstNumber"`
	PANNumber  string `json:"panNumber"`
	ColorCode  string `json:"colorCode"`
	Status     string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateBranchRequest) Validate() error {
	request.Code = strings.TrimSpace(request.Code)
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Address = strings.TrimSpace(request.Address)
	request.GSTNumber = strings.TrimSpace(request.GSTNumber)
	request.PANNumber = strings.TrimSpace(request.PANNumber)
	request.ColorCode = strings.TrimSpace(request.ColorCode)
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))

	if err := validateBranchFields(request.Code, request.Name, request.Email, request.Address, request.GSTNumber, request.PANNumber, request.ColorCode, request.Status); err != nil {
		return err
	}
	return nil
}

// UpdateBranchRequest is the body for PUT /branches/{id}. An empty status keeps
// the current one.
type UpdateBranchRequest struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	UserCityID *int64 `json:"userCityId"`
	CustomerID *int64 `json:"customerId"`
	GSTNumber  string `json:"gstNumber"`
	PANNumber  string `json:"panNumber"`
	ColorCode  string `json:"colorCode"`
	Status     string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateBranchRequest) Validate() error {
	request.Code = strings.TrimSpace(request.Code)
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Address = strings.TrimSpace(request.Address)
	request.GSTNumber = strings.TrimSpace(request.GSTNumber)
	request.PANNumber = strings.TrimSpace(request.PANNumber)
	request.ColorCode = strings.TrimSpace(request.ColorCode)
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))

	if err := validateBranchFields(request.Code, request.Name, request.Email, request.Address, request.GSTNumber, request.PANNumber, request.ColorCode, request.Status); err != nil {
		return err
	}
	return nil
}

// validateBranchFields checks the fields shared by create and update.
func validateBranchFields(code, name, email, address, gstNumber, panNumber, colorCode, status string) error {
	if !codePattern.MatchString(code) {
		return response.NewError(http.StatusBadRequest, "Please enter valid code")
	}
	if len(name) == 0 || len(name) > 200 {
		return response.NewError(http.StatusBadRequest, "Please enter valid name")
	}
	if email != "" {
		if address, err := mail.ParseAddress(email); err != nil || address.Address != email || len(email) > 255 {
			return response.NewError(http.StatusBadRequest, "Please enter valid email")
		}
	}
	if len(address) > 500 {
		return response.NewError(http.StatusBadRequest, "Please enter valid address")
	}
	if gstNumber != "" && !gstPattern.MatchString(gstNumber) {
		return response.NewError(http.StatusBadRequest, "Please enter valid GST number")
	}
	if panNumber != "" && !panPattern.MatchString(panNumber) {
		return response.NewError(http.StatusBadRequest, "Please enter valid PAN number")
	}
	if colorCode != "" && !colorCodePattern.MatchString(colorCode) {
		return response.NewError(http.StatusBadRequest, "Please enter valid color code")
	}
	if status != "" && status != BranchStatusActive && status != BranchStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}
