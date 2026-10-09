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
func (request *ItemRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	if request.Name == "" {
		return response.ValidName
	}
	if request.Price < 0 {
		return response.ValidPrice
	}
	return nil
}
