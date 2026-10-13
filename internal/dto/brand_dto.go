package dto

import "strings"

// CreateBrandRequest is the body for POST /brands.
type CreateBrandRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateBrandRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Description = strings.TrimSpace(request.Description)

	return validateCatalogFields(request.Name, request.Code, request.Description)
}

// UpdateBrandRequest is the body for PUT /brands/{id}.
type UpdateBrandRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateBrandRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Description = strings.TrimSpace(request.Description)

	return validateCatalogFields(request.Name, request.Code, request.Description)
}
