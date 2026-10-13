package dto

import (
	"strings"

	"7th-sense-server/internal/response"
)

// CreateControllerRequest is the body for POST /controllers.
type CreateControllerRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Variant     string `json:"variant"`
	BrandID     *int64 `json:"brandId"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateControllerRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Variant = strings.TrimSpace(request.Variant)
	request.Description = strings.TrimSpace(request.Description)

	return validateControllerFields(request.Name, request.Code, request.Variant, request.BrandID, request.Description)
}

// UpdateControllerRequest is the body for PUT /controllers/{id}.
type UpdateControllerRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Variant     string `json:"variant"`
	BrandID     *int64 `json:"brandId"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateControllerRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Variant = strings.TrimSpace(request.Variant)
	request.Description = strings.TrimSpace(request.Description)

	return validateControllerFields(request.Name, request.Code, request.Variant, request.BrandID, request.Description)
}

// validateControllerFields checks the fields shared by create and update.
func validateControllerFields(name, code, variant string, brandID *int64, description string) error {
	if err := validateCatalogFields(name, code, description); err != nil {
		return err
	}
	if len(variant) > 100 {
		return response.ValidVariant
	}
	if brandID != nil && *brandID <= 0 {
		return response.ValidBrand
	}
	return nil
}
