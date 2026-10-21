package dto

import (
	"net/http"
	"strings"

	"7th-sense-server/internal/response"
)

// Lengths of the downtime_reasons columns (migration 027).
const (
	downtimeReasonCodeMaxLength  = 10
	downtimeReasonMaxLength      = 50
	DowntimeReasonStatusActive   = "active"
	DowntimeReasonStatusInactive = "inactive"
)

// DowntimeReasonFields are the Add / Edit Downtime Reason form fields shared by create and update.
type DowntimeReasonFields struct {
	UserID     int64  `json:"userId"`     // customer the reason belongs to
	ReasonCode string `json:"reasonCode"` // optional
	Reason     string `json:"reason"`
}

// CreateDowntimeReasonRequest is the body for POST /downtime-reasons.
type CreateDowntimeReasonRequest struct {
	DowntimeReasonFields
}

// Validate trims the fields and checks they are usable.
func (request *CreateDowntimeReasonRequest) Validate() error {
	return request.DowntimeReasonFields.validate()
}

// UpdateDowntimeReasonRequest is the body for PUT /downtime-reasons/{id}. A blank status keeps the saved one.
type UpdateDowntimeReasonRequest struct {
	DowntimeReasonFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateDowntimeReasonRequest) Validate() error {
	if err := request.DowntimeReasonFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != DowntimeReasonStatusActive && request.Status != DowntimeReasonStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *DowntimeReasonFields) validate() error {
	fields.ReasonCode = strings.TrimSpace(fields.ReasonCode)
	fields.Reason = strings.TrimSpace(fields.Reason)

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case fields.ReasonCode != "" && !validText(fields.ReasonCode, downtimeReasonCodeMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid reason code")
	case !validText(fields.Reason, downtimeReasonMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid reason")
	}
	return nil
}
