package dto

import "strings"

// CreateCategoryRequest is the body for POST /categories.
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateCategoryRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Description = strings.TrimSpace(request.Description)

	return validateCatalogFields(request.Name, request.Code, request.Description)
}

// UpdateCategoryRequest is the body for PUT /categories/{id}.
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateCategoryRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Description = strings.TrimSpace(request.Description)

	return validateCatalogFields(request.Name, request.Code, request.Description)
}
