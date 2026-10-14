package dto

import (
	"strings"

	"7th-sense-server/internal/response"
)

// CreateControllerRequest is the body for POST /controllers. The variant (828D, 0i-F, …)
// is part of the name, so there is no separate code or variant.
type CreateControllerRequest struct {
	Name        string `json:"name"`
	BrandID     *int64 `json:"brandId"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateControllerRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	return validateControllerFields(request.Name, request.BrandID, request.Description)
}

// UpdateControllerRequest is the body for PUT /controllers/{id}.
type UpdateControllerRequest struct {
	Name        string `json:"name"`
	BrandID     *int64 `json:"brandId"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateControllerRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	return validateControllerFields(request.Name, request.BrandID, request.Description)
}

// validateControllerFields checks the fields shared by create and update. The brand
// is required: the machine form lists the controllers of the selected brand.
func validateControllerFields(name string, brandID *int64, description string) error {
	if brandID == nil {
		return response.SelectBrand
	}
	if *brandID <= 0 {
		return response.ValidBrand
	}
	return validateCatalogFields(name, "", description)
}
