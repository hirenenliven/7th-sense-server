package dto

import (
	"strings"

	"7th-sense-server/internal/response"
)

// ItemRequest is the body for creating or updating an item.
type ItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Validate trims the fields and checks they are usable.
func (r *ItemRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)

	if r.Name == "" {
		return response.ValidName
	}
	if r.Price < 0 {
		return response.ValidPrice
	}
	return nil
}
