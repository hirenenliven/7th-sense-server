package dto

import (
	"net/http"
	"strings"

	"7th-sense-server/internal/response"
)

// Lengths of the idle_reasons columns (migration 027).
const (
	idleReasonCodeMaxLength  = 10
	idleReasonMaxLength      = 50
	IdleReasonStatusActive   = "active"
	IdleReasonStatusInactive = "inactive"
)

// IdleReasonFields are the Add / Edit Idle Reason form fields shared by create and update.
type IdleReasonFields struct {
	UserID     int64  `json:"userId"`     // customer the reason belongs to
	ReasonCode string `json:"reasonCode"` // optional
	Reason     string `json:"reason"`
}

// CreateIdleReasonRequest is the body for POST /idle-reasons.
type CreateIdleReasonRequest struct {
	IdleReasonFields
}

// Validate trims the fields and checks they are usable.
func (request *CreateIdleReasonRequest) Validate() error {
	return request.IdleReasonFields.validate()
}

// UpdateIdleReasonRequest is the body for PUT /idle-reasons/{id}. A blank status keeps the saved one.
type UpdateIdleReasonRequest struct {
	IdleReasonFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateIdleReasonRequest) Validate() error {
	if err := request.IdleReasonFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != IdleReasonStatusActive && request.Status != IdleReasonStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *IdleReasonFields) validate() error {
	fields.ReasonCode = strings.TrimSpace(fields.ReasonCode)
	fields.Reason = strings.TrimSpace(fields.Reason)

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case fields.ReasonCode != "" && !validText(fields.ReasonCode, idleReasonCodeMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid reason code")
	case !validText(fields.Reason, idleReasonMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid reason")
	}
	return nil
}
