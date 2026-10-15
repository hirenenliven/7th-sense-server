package dto

import (
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"7th-sense-server/internal/response"
)

// Lengths of the clienteles columns (migrations 017, 018).
const (
	clientCodeMaxLength     = 50
	clienteleNameMaxLength  = 250
	clienteleEmailMaxLength = 255
	clienteleAddressMaxLen  = 250
	clienteleCityMaxLength  = 100
	clienteleRegionMaxLen   = 100 // country_name, state_name
	ClienteleStatusActive   = "active"
	ClienteleStatusInactive = "inactive"
)

var (
	internationalPhonePattern = regexp.MustCompile(`^\+[0-9]{7,15}$`) // +<code><number>
)

// ClienteleFields are the Add / Edit Clientele form fields shared by create and update.
type ClienteleFields struct {
	UserID        int64  `json:"userId"` // customer the client belongs to
	ClientCode    string `json:"clientCode"`
	Name          string `json:"name"`
	Email         string `json:"email"` // optional
	ContactNumber string `json:"contactNumber"`
	Address       string `json:"address"`
	CountryName   string `json:"countryName"` // e.g. "India"
	StateName     string `json:"stateName"`   // e.g. "Gujarat"
	City          string `json:"city"`
}

// CreateClienteleRequest is the body for POST /clienteles.
type CreateClienteleRequest struct {
	ClienteleFields
}

// Validate trims the fields and checks they are usable.
func (request *CreateClienteleRequest) Validate() error {
	return request.ClienteleFields.validate()
}

// UpdateClienteleRequest is the body for PUT /clienteles/{id}. A blank status keeps the saved one.
type UpdateClienteleRequest struct {
	ClienteleFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateClienteleRequest) Validate() error {
	if err := request.ClienteleFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != ClienteleStatusActive && request.Status != ClienteleStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *ClienteleFields) validate() error {
	fields.ClientCode = strings.TrimSpace(fields.ClientCode)
	fields.Name = strings.TrimSpace(fields.Name)
	fields.Email = strings.ToLower(strings.TrimSpace(fields.Email))
	fields.ContactNumber = strings.NewReplacer(" ", "", "-", "").Replace(strings.TrimSpace(fields.ContactNumber))
	fields.Address = strings.TrimSpace(fields.Address)
	fields.CountryName = strings.TrimSpace(fields.CountryName)
	fields.StateName = strings.TrimSpace(fields.StateName)
	fields.City = strings.TrimSpace(fields.City)

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case !validText(fields.ClientCode, clientCodeMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid client ID")
	case !validText(fields.Name, clienteleNameMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid client name")
	case fields.Email != "" && !validClienteleEmail(fields.Email):
		return response.NewError(http.StatusBadRequest, "Please enter valid client email")
	case !internationalPhonePattern.MatchString(fields.ContactNumber):
		return response.NewError(http.StatusBadRequest, "Please enter valid contact number")
	case !validText(fields.Address, clienteleAddressMaxLen):
		return response.NewError(http.StatusBadRequest, "Please enter valid address")
	case !validText(fields.CountryName, clienteleRegionMaxLen):
		return response.NewError(http.StatusBadRequest, "Please select valid country")
	case !validText(fields.StateName, clienteleRegionMaxLen):
		return response.NewError(http.StatusBadRequest, "Please select valid state")
	case !validText(fields.City, clienteleCityMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid city")
	}
	return nil
}

func validClienteleEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email && len(email) <= clienteleEmailMaxLength
}
