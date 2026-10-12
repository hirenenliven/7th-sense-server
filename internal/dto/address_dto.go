package dto

import (
	"regexp"
	"strings"

	"7th-sense-server/internal/response"
)

var (
	pinCodePattern = regexp.MustCompile(`^[0-9]{3,10}$`)
	phonePattern   = regexp.MustCompile(`^\+?[0-9]{7,15}$`)
)

// Address types accepted by CreateAddressRequest and UpdateAddressRequest.
const (
	AddressTypeHome     = "home"
	AddressTypeWork     = "work"
	AddressTypeShipping = "shipping"
	AddressTypeBilling  = "billing"
)

// CreateAddressRequest is the body for POST /addresses.
type CreateAddressRequest struct {
	UserID        int64  `json:"userId"`
	AddressType   string `json:"addressType"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	StateID       *int64 `json:"stateId"`
	CountryID     *int64 `json:"countryId"`
	PinCode       string `json:"pinCode"`
	ContactPerson string `json:"contactPerson"`
	ContactPhone  string `json:"contactPhone"`
	IsDefault     bool   `json:"isDefault"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateAddressRequest) Validate() error {
	request.AddressType = strings.ToLower(strings.TrimSpace(request.AddressType))
	request.AddressLine1 = strings.TrimSpace(request.AddressLine1)
	request.AddressLine2 = strings.TrimSpace(request.AddressLine2)
	request.City = strings.TrimSpace(request.City)
	request.PinCode = strings.TrimSpace(request.PinCode)
	request.ContactPerson = strings.TrimSpace(request.ContactPerson)
	request.ContactPhone = strings.TrimSpace(request.ContactPhone)

	if request.UserID <= 0 {
		return response.InvalidID
	}
	if request.AddressType == "" {
		request.AddressType = AddressTypeHome
	}
	if request.AddressType != AddressTypeHome && request.AddressType != AddressTypeWork &&
		request.AddressType != AddressTypeShipping && request.AddressType != AddressTypeBilling {
		return response.ValidAddress
	}
	if request.AddressLine1 == "" {
		return response.ValidAddress
	}
	if request.City == "" {
		return response.ValidCity
	}
	if request.PinCode != "" && !pinCodePattern.MatchString(request.PinCode) {
		return response.ValidAddress
	}
	if request.ContactPhone != "" && !phonePattern.MatchString(request.ContactPhone) {
		return response.ValidMobile
	}
	return nil
}

// UpdateAddressRequest is the body for PUT /addresses/{id}.
type UpdateAddressRequest struct {
	AddressType   string `json:"addressType"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	StateID       *int64 `json:"stateId"`
	CountryID     *int64 `json:"countryId"`
	PinCode       string `json:"pinCode"`
	ContactPerson string `json:"contactPerson"`
	ContactPhone  string `json:"contactPhone"`
	IsDefault     bool   `json:"isDefault"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateAddressRequest) Validate() error {
	request.AddressType = strings.ToLower(strings.TrimSpace(request.AddressType))
	request.AddressLine1 = strings.TrimSpace(request.AddressLine1)
	request.AddressLine2 = strings.TrimSpace(request.AddressLine2)
	request.City = strings.TrimSpace(request.City)
	request.PinCode = strings.TrimSpace(request.PinCode)
	request.ContactPerson = strings.TrimSpace(request.ContactPerson)
	request.ContactPhone = strings.TrimSpace(request.ContactPhone)

	if request.AddressType != "" && request.AddressType != AddressTypeHome && request.AddressType != AddressTypeWork &&
		request.AddressType != AddressTypeShipping && request.AddressType != AddressTypeBilling {
		return response.ValidAddress
	}
	if request.AddressLine1 != "" && len(request.AddressLine1) > 255 {
		return response.ValidAddress
	}
	if request.AddressLine2 != "" && len(request.AddressLine2) > 255 {
		return response.ValidAddress
	}
	if request.City != "" && len(request.City) > 100 {
		return response.ValidCity
	}
	if request.PinCode != "" && !pinCodePattern.MatchString(request.PinCode) {
		return response.ValidAddress
	}
	if request.ContactPerson != "" && len(request.ContactPerson) > 100 {
		return response.ValidName
	}
	if request.ContactPhone != "" && !phonePattern.MatchString(request.ContactPhone) {
		return response.ValidMobile
	}
	return nil
}
