package dto

import (
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata" // time zone database for time.LoadLocation, so Windows hosts can validate time zones

	"7th-sense-server/internal/model"
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

// dayStartLayouts are the accepted formats of CustomerDetails.DayStartFrom.
var dayStartLayouts = []string{"15:04", "15:04:05"}

// CreateUserRequest is the body for POST /users: the login credentials plus
// the customer details.
type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	RoleID    int64  `json:"roleId"`
	CompanyID *int64 `json:"companyId"`
	CustomerDetails
}

// Validate trims the fields and checks they are usable.
func (request *CreateUserRequest) Validate() error {
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)

	if err := validateProfile(request.Username, request.Email, request.FirstName, request.LastName); err != nil {
		return err
	}
	// bcrypt only uses the first 72 bytes of a password.
	if len(request.Password) < 8 || len(request.Password) > 72 {
		return response.NewError(http.StatusBadRequest, "Please enter valid password")
	}
	// Role is mandatory
	if request.RoleID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid role")
	}
	return request.CustomerDetails.validate()
}

// UpdateUserRequest is the body for PUT /users/{id}. It replaces every field
// except the password. An empty status keeps the current one.
type UpdateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	RoleID    int64  `json:"roleId"`
	CompanyID *int64 `json:"companyId"`
	Status    string `json:"status"`
	CustomerDetails
}

// Validate trims the fields and checks they are usable.
func (request *UpdateUserRequest) Validate() error {
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.FirstName = strings.TrimSpace(request.FirstName)
	request.LastName = strings.TrimSpace(request.LastName)
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))

	if err := validateProfile(request.Username, request.Email, request.FirstName, request.LastName); err != nil {
		return err
	}
	if request.Status != "" && request.Status != UserStatusActive && request.Status != UserStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	// Role is mandatory
	if request.RoleID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid role")
	}
	return request.CustomerDetails.validate()
}

// CustomerDetails holds the customer form fields shared by create and update.
// They are saved in users (branch, primary phone), company_profiles (company
// details), addresses (address) and user_profiles (logo, contact, settings).
type CustomerDetails struct {
	// Company details
	BranchID    *int64 `json:"branchId"`
	CompanyName string `json:"companyName"`
	FinixID     *int64 `json:"finixId"`
	GSTNumber   string `json:"gstNumber"`
	PANNumber   string `json:"panNumber"`
	Website     string `json:"website"`
	CompanyLogo string `json:"companyLogo"` // path or URL of the uploaded logo

	// Address
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	CountryID    *int64 `json:"countryId"`
	StateID      *int64 `json:"stateId"`
	City         string `json:"city"`
	ZipCode      string `json:"zipCode"`

	// Contact & settings
	ContactPerson  string `json:"contactPerson"`
	PrimaryPhone   string `json:"primaryPhone"`
	SecondaryPhone string `json:"secondaryPhone"`
	TimeZone       string `json:"timeZone"`     // IANA name, e.g. "Asia/Kolkata"
	DayStartFrom   string `json:"dayStartFrom"` // "HH:MM", e.g. "06:00"
}

// validate trims the customer fields and checks they are usable.
func (details *CustomerDetails) validate() error {
	details.CompanyName = strings.TrimSpace(details.CompanyName)
	details.GSTNumber = strings.ToUpper(strings.TrimSpace(details.GSTNumber))
	details.PANNumber = strings.ToUpper(strings.TrimSpace(details.PANNumber))
	details.Website = strings.TrimSpace(details.Website)
	details.CompanyLogo = strings.TrimSpace(details.CompanyLogo)
	details.AddressLine1 = strings.TrimSpace(details.AddressLine1)
	details.AddressLine2 = strings.TrimSpace(details.AddressLine2)
	details.City = strings.TrimSpace(details.City)
	details.ZipCode = strings.TrimSpace(details.ZipCode)
	details.ContactPerson = strings.TrimSpace(details.ContactPerson)
	details.PrimaryPhone = strings.TrimSpace(details.PrimaryPhone)
	details.SecondaryPhone = strings.TrimSpace(details.SecondaryPhone)
	details.TimeZone = strings.TrimSpace(details.TimeZone)
	details.DayStartFrom = strings.TrimSpace(details.DayStartFrom)

	// Company details
	if details.BranchID == nil || *details.BranchID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid branch")
	}
	if details.CompanyName == "" || len(details.CompanyName) > 150 {
		return response.NewError(http.StatusBadRequest, "Please enter valid company name")
	}
	if details.FinixID != nil && *details.FinixID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid Finix ID")
	}
	if !gstPattern.MatchString(details.GSTNumber) {
		return response.NewError(http.StatusBadRequest, "Please enter valid GST number")
	}
	if !panPattern.MatchString(details.PANNumber) {
		return response.NewError(http.StatusBadRequest, "Please enter valid PAN number")
	}
	if details.Website != "" && !isWebsite(details.Website) {
		return response.NewError(http.StatusBadRequest, "Please enter valid website")
	}
	if len(details.CompanyLogo) > 255 {
		return response.NewError(http.StatusBadRequest, "Please enter valid image")
	}

	// Address
	if details.AddressLine1 == "" || len(details.AddressLine1) > 255 || len(details.AddressLine2) > 255 {
		return response.NewError(http.StatusBadRequest, "Please enter valid address")
	}
	if details.CountryID == nil || *details.CountryID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid country")
	}
	if details.StateID == nil || *details.StateID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please enter valid state")
	}
	if details.City == "" || len(details.City) > 100 {
		return response.NewError(http.StatusBadRequest, "Please enter valid city")
	}
	if !pinCodePattern.MatchString(details.ZipCode) {
		return response.NewError(http.StatusBadRequest, "Please enter valid zip code")
	}

	// Contact & settings
	if details.ContactPerson == "" || len(details.ContactPerson) > 100 {
		return response.NewError(http.StatusBadRequest, "Please enter valid contact person")
	}
	if !mobilePattern.MatchString(details.PrimaryPhone) {
		return response.NewError(http.StatusBadRequest, "Please enter valid mobile")
	}
	if details.SecondaryPhone != "" && !phonePattern.MatchString(details.SecondaryPhone) {
		return response.NewError(http.StatusBadRequest, "Please enter valid mobile")
	}
	if !isTimeZone(details.TimeZone) {
		return response.NewError(http.StatusBadRequest, "Please enter valid time zone")
	}
	if _, err := ParseDayStart(details.DayStartFrom); err != nil {
		return response.NewError(http.StatusBadRequest, "Please enter valid day start time")
	}
	return nil
}

// ParseDayStart parses a day start time such as "06:00" or "06:00:00".
func ParseDayStart(value string) (time.Time, error) {
	var parseErr error
	for _, layout := range dayStartLayouts {
		dayStart, err := time.Parse(layout, value)
		if err == nil {
			return dayStart, nil
		}
		parseErr = err
	}
	return time.Time{}, parseErr
}

// isWebsite reports whether website is an http(s) URL with a host. The scheme
// may be left out, as in "www.example.com".
func isWebsite(website string) bool {
	if len(website) > 255 {
		return false
	}
	if !strings.Contains(website, "://") {
		website = "https://" + website
	}
	parsed, err := url.ParseRequestURI(website)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && strings.Contains(parsed.Host, ".")
}

// isTimeZone reports whether timeZone is an IANA time zone name such as "Asia/Kolkata".
func isTimeZone(timeZone string) bool {
	// LoadLocation accepts "" (UTC) and "Local" (the server's zone), which are not real names.
	if timeZone == "" || timeZone == "Local" || len(timeZone) > 50 {
		return false
	}
	_, err := time.LoadLocation(timeZone)
	return err == nil
}

// validateProfile checks the user fields shared by create and update.
func validateProfile(username, email, firstName, lastName string) error {
	if !usernamePattern.MatchString(username) {
		return response.NewError(http.StatusBadRequest, "Please enter valid username")
	}
	// ParseAddress also accepts "Name <a@b.com>", so require a bare address.
	if address, err := mail.ParseAddress(email); err != nil || address.Address != email || len(email) > 255 {
		return response.NewError(http.StatusBadRequest, "Please enter valid email")
	}
	if len(firstName) > 100 || len(lastName) > 100 {
		return response.NewError(http.StatusBadRequest, "Please enter valid name")
	}
	return nil
}

// UserResponse is a user with its customer details. A part the user does not
// have yet (for example a user created before the customer form) is null.
type UserResponse struct {
	model.User
	Profile        *model.UserProfile    `json:"profile"`
	CompanyProfile *model.CompanyProfile `json:"companyProfile"`
	Address        *model.Address        `json:"address"`
}
