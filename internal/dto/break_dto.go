package dto

import (
	"net/http"
	"strings"

	"7th-sense-server/internal/response"
)

// Break statuses accepted by UpdateBreakRequest.
const (
	BreakStatusActive   = "active"
	BreakStatusInactive = "inactive"
)

const breakNameMaxLength = 100

// BreakFields are the Manage Breaks form fields shared by create and update.
type BreakFields struct {
	UserID    int64  `json:"userId"`  // customer the break belongs to
	ShiftID   *int64 `json:"shiftId"` // one of the customer's shifts; nil = no shift yet
	Name      string `json:"name"`
	StartTime string `json:"startTime"` // "HH:MM"
	EndTime   string `json:"endTime"`   // "HH:MM"; not after start = past midnight
}

// CreateBreakRequest is the body for POST /breaks.
type CreateBreakRequest struct {
	BreakFields
}

// Validate trims the fields and checks they are usable.
func (request *CreateBreakRequest) Validate() error {
	return request.BreakFields.validate()
}

// UpdateBreakRequest is the body for PUT /breaks/{id}. A blank status keeps the saved one.
type UpdateBreakRequest struct {
	BreakFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateBreakRequest) Validate() error {
	if err := request.BreakFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != BreakStatusActive && request.Status != BreakStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *BreakFields) validate() error {
	fields.Name = strings.TrimSpace(fields.Name)
	fields.StartTime = strings.TrimSpace(fields.StartTime)
	fields.EndTime = strings.TrimSpace(fields.EndTime)

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case fields.ShiftID != nil && *fields.ShiftID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select valid shift")
	case !validText(fields.Name, breakNameMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid break name")
	case !clockPattern.MatchString(fields.StartTime):
		return response.NewError(http.StatusBadRequest, "Please enter valid break start time")
	case !clockPattern.MatchString(fields.EndTime):
		return response.NewError(http.StatusBadRequest, "Please enter valid break end time")
	case fields.StartTime == fields.EndTime:
		return response.NewError(http.StatusBadRequest, "Break end time must differ from the start time")
	}
	return nil
}
