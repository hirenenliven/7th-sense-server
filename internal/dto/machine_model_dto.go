package dto

import (
	"strings"

	"7th-sense-server/internal/response"
)

// MachineModelRequest is the body for POST /models and PUT /models/{id}.
type MachineModelRequest struct {
	ControllerID int64  `json:"controllerId"`
	Name         string `json:"name"`
}

// Validate trims the name and checks the fields are usable.
func (request *MachineModelRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)

	if request.ControllerID <= 0 {
		return response.SelectController
	}
	if request.Name == "" || len(request.Name) > maxCatalogNameLength {
		return response.ValidName
	}
	return nil
}
