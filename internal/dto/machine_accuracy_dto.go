package dto

import (
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// MachineAccuracyNameMaxLength is the size of machine_accuracies.name.
const MachineAccuracyNameMaxLength = 50

// MachineAccuracyRequest is the body for POST /machine-accuracies and PUT
// /machine-accuracies/{id}: the accuracy's name.
type MachineAccuracyRequest struct {
	Name string `json:"name"`
}

// Validate trims the name, joins runs of spaces and checks its length.
func (request *MachineAccuracyRequest) Validate() error {
	request.Name = strings.Join(strings.Fields(request.Name), " ")
	if request.Name == "" {
		return response.EnterMachineAccuracy
	}
	if utf8.RuneCountInString(request.Name) > MachineAccuracyNameMaxLength {
		return response.MachineAccuracyTooLong
	}
	return nil
}
