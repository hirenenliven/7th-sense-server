package dto

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// MaintenanceTaskDescriptionMaxLength is the size of maintenance_tasks.description.
const MaintenanceTaskDescriptionMaxLength = 100

// Maintenance task statuses (maintenance_tasks.status).
const (
	MaintenanceTaskStatusActive   = "active"
	MaintenanceTaskStatusInactive = "inactive"
)

// CreateMaintenanceTaskRequest is the body for POST /maintenance-tasks: a new task for the master
// list, added from the Add Operator Task form.
type CreateMaintenanceTaskRequest struct {
	Description string `json:"description"`
}

// Validate trims the description, joins runs of spaces and checks its length.
func (request *CreateMaintenanceTaskRequest) Validate() error {
	request.Description = strings.Join(strings.Fields(request.Description), " ")
	if request.Description == "" {
		return response.EnterTask
	}
	if utf8.RuneCountInString(request.Description) > MaintenanceTaskDescriptionMaxLength {
		return response.TaskTooLong
	}
	return nil
}

// UpdateMaintenanceTaskRequest is the body for PUT /maintenance-tasks/{id}: the description and,
// optionally, the status (blank keeps it). An inactive task is left out of the Task dropdown; the
// operator tasks that already use it keep it.
type UpdateMaintenanceTaskRequest struct {
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Validate checks the description as on create, and the status.
func (request *UpdateMaintenanceTaskRequest) Validate() error {
	descriptionRequest := CreateMaintenanceTaskRequest{Description: request.Description}
	if err := descriptionRequest.Validate(); err != nil {
		return err
	}
	request.Description = descriptionRequest.Description
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != MaintenanceTaskStatusActive && request.Status != MaintenanceTaskStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

// OperatorTaskFields are the fields of the Add / Edit Operator Task form: the customer, the
// machine type (a catalog model; nil = every machine of the customer) and the task.
type OperatorTaskFields struct {
	UserID            int64  `json:"userId"`
	ModelID           *int64 `json:"modelId"`
	MaintenanceTaskID int64  `json:"maintenanceTaskId"`
}

// Validate checks the ids are given (and the model, when given, is a positive id).
func (fields *OperatorTaskFields) Validate() error {
	switch {
	case fields.UserID <= 0:
		return response.SelectCustomer
	case fields.ModelID != nil && *fields.ModelID <= 0:
		return response.SelectValidMachineType
	case fields.MaintenanceTaskID <= 0:
		return response.SelectTask
	}
	return nil
}

// CreateOperatorTaskRequest is the body for POST /operator-tasks.
type CreateOperatorTaskRequest struct {
	OperatorTaskFields
}

// UpdateOperatorTaskRequest is the body for PUT /operator-tasks/{id}.
type UpdateOperatorTaskRequest struct {
	OperatorTaskFields
}
