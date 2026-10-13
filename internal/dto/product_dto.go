package dto

import (
	"strings"

	"7th-sense-server/internal/response"
)

// CreateProductRequest is the body for POST /products. An empty image gets
// the default product image.
type CreateProductRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CategoryID  *int64 `json:"categoryId"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateProductRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Image = strings.TrimSpace(request.Image)
	request.Description = strings.TrimSpace(request.Description)

	return validateProductFields(request.Name, request.Code, request.CategoryID, request.Image, request.Description)
}

// UpdateProductRequest is the body for PUT /products/{id}. An empty image gets
// the default product image.
type UpdateProductRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CategoryID  *int64 `json:"categoryId"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateProductRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Code = strings.ToUpper(strings.TrimSpace(request.Code))
	request.Image = strings.TrimSpace(request.Image)
	request.Description = strings.TrimSpace(request.Description)

	return validateProductFields(request.Name, request.Code, request.CategoryID, request.Image, request.Description)
}

// validateProductFields checks the fields shared by create and update.
func validateProductFields(name, code string, categoryID *int64, image, description string) error {
	if err := validateCatalogFields(name, code, description); err != nil {
		return err
	}
	if categoryID != nil && *categoryID <= 0 {
		return response.ValidCategory
	}
	if len(image) > 255 {
		return response.ValidImage
	}
	return nil
}
